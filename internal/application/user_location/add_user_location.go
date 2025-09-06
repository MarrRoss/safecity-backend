package user_location

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/adapter/db_storage/response"
	"awesomeProjectDDD/internal/adapter/telegram"
	"awesomeProjectDDD/internal/application"
	"awesomeProjectDDD/internal/domain/model/entity"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
	"fmt"
	"time"
)

type AddUserLocationCommand struct {
	UserID    string
	Latitude  float64
	Longitude float64
	Battery   int
}

type AddUserLocationHandler struct {
	locationStorage            db.UserLocationRepository
	userStorage                db.UserRepository
	notificationLogStorage     db.NotificationLogRepository
	familyStorage              db.FamilyRepository
	membershipStorage          db.FamilyMembershipRepository
	zoneStorage                db.ZoneRepository
	notificationSettingStorage db.NotificationSettingRepository
	frequencyStorage           db.NotificationFrequencyRepository
	observer                   *observability.Observability
	telegramService            *telegram.Service
}

func NewAddUserLocationHandler(
	locationStorage db.UserLocationRepository,
	userStorage db.UserRepository,
	notificationLogStorage db.NotificationLogRepository,
	familyStorage db.FamilyRepository,
	membershipStorage db.FamilyMembershipRepository,
	zoneStorage db.ZoneRepository,
	notificationSettingStorage db.NotificationSettingRepository,
	frequencyStorage db.NotificationFrequencyRepository,
	observer *observability.Observability,
	telegramService *telegram.Service,
) *AddUserLocationHandler {
	return &AddUserLocationHandler{
		userStorage:                userStorage,
		locationStorage:            locationStorage,
		notificationLogStorage:     notificationLogStorage,
		familyStorage:              familyStorage,
		membershipStorage:          membershipStorage,
		zoneStorage:                zoneStorage,
		notificationSettingStorage: notificationSettingStorage,
		frequencyStorage:           frequencyStorage,
		observer:                   observer,
		telegramService:            telegramService,
	}
}

func (h *AddUserLocationHandler) Handle(ctx context.Context, cmd AddUserLocationCommand) (value_object.ID, []value_object.ID, error) {
	extID, err := value_object.NewIDFromString(cmd.UserID)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid user id")
		return value_object.ID{}, []value_object.ID{}, err
	}
	user, err := h.userStorage.GetUserByExternalID(ctx, extID)
	if err != nil {
		if errors.Is(err, adapter.ErrUserNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("user not found")
			return value_object.ID{}, []value_object.ID{}, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting user")
			return value_object.ID{}, []value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting user")
		return value_object.ID{}, []value_object.ID{}, err
	}
	if !user.Tracking {
		h.observer.Logger.Error().Msg("user is not tracked")
		return value_object.ID{}, []value_object.ID{}, application.ErrUserNotTracked
	}

	lat, err := value_object.NewLatitude(cmd.Latitude)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to add latitude")
		return value_object.ID{}, []value_object.ID{}, fmt.Errorf("failed to add latitude: %w", err)
	}
	lng, err := value_object.NewLongitude(cmd.Longitude)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to add longitude")
		return value_object.ID{}, []value_object.ID{}, fmt.Errorf("failed to add longitude: %w", err)
	}
	location := value_object.NewLocation(lat, lng)
	battery, err := value_object.NewBatteryThreshold(cmd.Battery)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid battery value")
		return value_object.ID{}, nil, fmt.Errorf("invalid battery: %w", err)
	}
	userLocation := entity.NewUserLocation(user.ID, location, battery)
	err = h.locationStorage.AddUserLocation(ctx, userLocation)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while adding user location")
			return value_object.ID{}, []value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while adding user location")
		return value_object.ID{}, []value_object.ID{}, err
	}
	//fmt.Printf("user location added: %v\n", userLocation.ID)

	userFamilies, err := h.membershipStorage.GetFamiliesByUserID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting user families")
			return value_object.ID{}, []value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting user families")
		return value_object.ID{}, []value_object.ID{}, err
	}
	if len(userFamilies) == 0 {
		h.observer.Logger.Info().Msg("no families found for user")
		return userLocation.ID, []value_object.ID{}, nil
	}
	familyID := userFamilies[0].ID
	//h.observer.Logger.Println("family id: ", familyID)
	zones, err := h.zoneStorage.CoordinatesInFamilyZones(ctx, familyID, cmd.Longitude, cmd.Latitude)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting family zone")
			return value_object.ID{}, []value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting family zone")
		return value_object.ID{}, []value_object.ID{}, err
	}
	var currZone *entity.Zone
	if len(zones) > 0 {
		currZone = zones[0]
		fmt.Printf("curr zone id: %v\n", currZone.ID)
	}
	prevCtx, err := h.locationStorage.GetLastZoneContext(ctx, user.ID)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting last location context")
			return value_object.ID{}, []value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting last location context")
		return value_object.ID{}, []value_object.ID{}, err
	}
	var prevZoneID *value_object.ID
	if prevCtx != nil && prevCtx.ZoneID != nil && prevCtx.NotifyType != "exit" {
		id, err := value_object.NewIDFromString(*prevCtx.ZoneID)
		if err != nil {
			return value_object.ID{}, []value_object.ID{}, err
		}
		prevZoneID = &id
	}
	var currZoneID *value_object.ID
	if currZone != nil {
		currZoneID = &currZone.ID
	}

	notifyType := determineContext(prevZoneID, currZoneID)
	//fmt.Printf("notifyType: %v\n", notifyType)
	var zoneForCtx *value_object.ID
	switch notifyType {
	case "entry", "inside":
		zoneForCtx = currZoneID
	case "exit":
		zoneForCtx = prevZoneID
	case "none":
		zoneForCtx = nil
	}

	ctxType := notifyType
	ctxZone := zoneForCtx
	// если дублирующий выход или вход — вместо него сохраняем none/NULL
	if (notifyType == "exit" || notifyType == "entry") &&
		prevCtx != nil &&
		prevCtx.NotifyType == notifyType &&
		prevCtx.ZoneID != nil &&
		zoneForCtx != nil &&
		*prevCtx.ZoneID == zoneForCtx.String() {
		ctxType = "none"
		ctxZone = nil
	}
	//fmt.Printf("ctxType: %v\n", ctxType)
	err = h.locationStorage.AddLocationContext(ctx, userLocation.ID, ctxZone, ctxType, nil)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while adding location context")
			return value_object.ID{}, []value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while adding location context")
		return value_object.ID{}, []value_object.ID{}, err
	}
	//fmt.Printf("location context added: %v\n", userLocation.ID)

	batteryNotifyType := "battery"
	if err := h.locationStorage.AddLocationContext(ctx, userLocation.ID, nil, batteryNotifyType, &cmd.Battery); err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while adding location context")
			return value_object.ID{}, []value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while adding location context")
		return value_object.ID{}, []value_object.ID{}, err
	}
	//fmt.Printf("battery context added: %v\n", userLocation.ID)

	var notificationLogsIDs []value_object.ID
	switch ctxType {
	case "entry", "exit":
		fmt.Printf("entry or exit case")
		if ctxZone != nil {
			pushIDs, err := h.CreateNotificationLogsByZone(ctx, user, *ctxZone, ctxType)
			if err != nil {
				h.observer.Logger.Error().Err(err).Msg("failed to create zone notification logs")
				return userLocation.ID, nil, fmt.Errorf("failed to create zone notification logs: %w", err)
			}
			if len(pushIDs) > 0 {
				notificationLogsIDs = append(notificationLogsIDs, pushIDs...)
				//fmt.Printf("zone notification logs created: %v\n", notificationLogsIDs)
			}
			h.observer.Logger.Print("case entry/exit")
		}
	case "inside":
		fmt.Printf("inside entry")
		//if ctxZone != nil {
		//	fmt.Printf("context zone: %v\n", *ctxZone)
		//}
		//if currZone != nil {
		//	fmt.Printf("current zone: %v\n", *currZone)
		//}
		if ctxZone != nil && currZone != nil && currZone.Safety.Bool() == false {
			fmt.Printf("inside case")
			pushIDs, err := h.PushInsideNotifications(ctx, user, *ctxZone, ctxType)
			if err != nil {
				h.observer.Logger.Error().Err(err).Msg("failed to create zone notification logs")
				return userLocation.ID, nil, fmt.Errorf("failed to create zone notification logs: %w", err)
			}
			if len(pushIDs) > 0 {
				notificationLogsIDs = append(notificationLogsIDs, pushIDs...)
				fmt.Printf("zone notification logs created: %v\n", notificationLogsIDs)
			}
			h.observer.Logger.Print("case inside")
		}
	case "none":
		h.observer.Logger.Info().Msg("no push notifications")
		//return userLocation.ID, []value_object.ID{}, nil
	}

	if cmd.Battery > 0 {
		fmt.Printf("battery: %v\n", cmd.Battery)
		batteryPushIDs, err := h.PushBatteryNotifications(ctx, user, cmd.Battery)
		if err != nil {
			h.observer.Logger.Error().Err(err).Msg("failed to create battery notification logs")
			return userLocation.ID, nil, fmt.Errorf("failed to create battery notification logs: %w", err)
		}
		if len(batteryPushIDs) > 0 {
			notificationLogsIDs = append(notificationLogsIDs, batteryPushIDs...)
			fmt.Printf("battery notification logs created: %v\n", notificationLogsIDs)
		}
		h.observer.Logger.Print("case battery")
	}
	return userLocation.ID, notificationLogsIDs, nil
}

func determineContext(prev, curr *value_object.ID) string {
	if prev == nil && curr == nil {
		return "none"
	}
	if prev == nil {
		return "entry"
	}
	if curr == nil {
		return "exit"
	}
	if *prev == *curr {
		return "inside"
	}
	return "none"
}

func (h *AddUserLocationHandler) PushBatteryNotifications(
	ctx context.Context,
	child *entity.User,
	battery int,
) ([]value_object.ID, error) {
	settings, err := h.notificationSettingStorage.GetBatteryNotificationSettingsByChild(ctx, child.ID)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting battery notification settings")
			return []value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting battery notification settings")
		return []value_object.ID{}, err
	}
	if len(settings) == 0 {
		h.observer.Logger.Info().Msg("no notification settings found for user battery")
		return []value_object.ID{}, nil
	}
	fmt.Printf("battery settings: %v\n", settings)

	var toCheck []string
	thresholdMap := make(map[string]int, len(settings))
	freqMap := make(map[string]time.Duration, len(settings))
	for _, s := range settings {
		if battery <= s.Threshold {
			toCheck = append(toCheck, s.NotifySettingID)
			thresholdMap[s.NotifySettingID] = s.Threshold
			freqMap[s.NotifySettingID] = s.Frequency
		}
	}
	if len(toCheck) == 0 {
		h.observer.Logger.Info().Msg("battery is ok")
		return []value_object.ID{}, nil
	}
	fmt.Printf("toCheck: %v\n", toCheck)

	var settingVOs []value_object.ID
	for _, id := range toCheck {
		vid, err := value_object.NewIDFromString(id)
		if err != nil {
			h.observer.Logger.Trace().Err(err).Msg("invalid notification setting id")
			return []value_object.ID{}, err
		}
		settingVOs = append(settingVOs, vid)
	}
	lastLogs, err := h.notificationLogStorage.GetLastLogsByNotificationIDs(ctx, settingVOs)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting notification logs")
			return []value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting notification logs")
		return []value_object.ID{}, err
	}
	fmt.Printf("battery lastLogs: %v\n", lastLogs)

	lastLogMap := make(map[string]*response.NotificationLogDB, len(lastLogs))
	for _, l := range lastLogs {
		lastLogMap[l.NotificationID] = l
	}

	receivers, err := h.notificationSettingStorage.GetIntegratedReceiversBySettingIDs(ctx, settingVOs)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting notification receivers")
			return []value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting notification receivers")
		return []value_object.ID{}, err
	}
	if len(receivers) == 0 {
		h.observer.Logger.Info().Msg("no integrated receivers found for user location")
		return []value_object.ID{}, nil
	}
	fmt.Printf("battery receivers: %v\n", receivers)

	now := time.Now().UTC()
	var pushed []value_object.ID

	for _, rcv := range receivers {
		sid := rcv.NotifySettingID.String()
		if rcv.System.ID != 1 {
			continue
		}
		freq, ok := freqMap[sid]
		if !ok {
			continue
		}
		var elapsed time.Duration
		if last := lastLogMap[sid]; last != nil {
			elapsed = now.Sub(last.CreatedAt)
			fmt.Printf("battery elapsed: %v\n", elapsed)
		}
		if lastLogMap[sid] != nil && elapsed < freq {
			fmt.Printf("skipping battery notification: %v\n", elapsed)
			continue
		}
		//last := lastLogMap[sid]
		//if last != nil && now.Sub(last.CreatedAt) < freq {
		//	continue
		//}
		if err := h.telegramService.SendMessage(ctx, rcv.System.ExternalID,
			fmt.Sprintf("Уровень батареи ребёнка %s: %d%%", child.Name.FirstName.String(), battery),
		); err != nil {
			h.observer.Logger.Trace().Err(err).Msg("failed to send battery notification")
			continue
		}
		fmt.Printf("send message")

		settingID, err := value_object.NewIDFromString(sid)
		if err != nil {
			h.observer.Logger.Trace().Err(err).Msg("invalid notification setting id")
			continue
		}
		log, err := entity.NewNotificationLog(settingID, "battery", 1)
		if err != nil {
			h.observer.Logger.Trace().Err(err).Msg("failed to create notification log")
			continue
		}
		fmt.Printf("battery notify log: %v\n", log)
		if err := h.notificationLogStorage.AddNotificationLog(ctx, log); err != nil {
			if errors.Is(err, adapter.ErrStorage) {
				h.observer.Logger.Error().Err(err).Msg("database error while adding notification log")
				continue
			}
			h.observer.Logger.Error().Err(err).Msg("unexpected error while adding notification log")
			continue
		}
		pushed = append(pushed, log.ID)
		fmt.Printf("notify log: %v\n", log)
	}
	return pushed, nil
}

func (h *AddUserLocationHandler) CreateNotificationLogsByZone(
	ctx context.Context,
	child *entity.User,
	zoneID value_object.ID,
	notifyType string,
) ([]value_object.ID, error) {
	zone, err := h.zoneStorage.GetZone(ctx, zoneID)
	if err != nil {
		if errors.Is(err, adapter.ErrZoneNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("zone not found")
			return []value_object.ID{}, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting zone")
			return []value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting zone")
		return []value_object.ID{}, err
	}
	settingIDs, err := h.notificationSettingStorage.GetDangerZoneNotificationSettings(ctx, child.ID, zoneID)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting zone notification settings")
			return []value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting zone notification settings")
		return []value_object.ID{}, err
	}
	if len(settingIDs) == 0 {
		h.observer.Logger.Info().Msg("no notification settings found for user location")
		return []value_object.ID{}, nil
	}
	//fmt.Printf("settingIDs: %v\n", settingIDs)
	receivers, err := h.notificationSettingStorage.GetIntegratedReceiversBySettingIDs(ctx, settingIDs)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting notification receivers")
			return []value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting notification receivers")
		return []value_object.ID{}, err
	}
	if len(receivers) == 0 {
		h.observer.Logger.Info().Msg("no integrated receivers found for user location")
		return []value_object.ID{}, nil
	}
	//fmt.Printf("receivers: %v\n", receivers)
	var notifyLogsIDs []value_object.ID
	for _, receiver := range receivers {
		if receiver.System.ID != 1 {
			continue
		}
		var contextType string
		if notifyType == "entry" {
			contextType = "вошел в"
		}
		if notifyType == "exit" {
			contextType = "вышел из"
		}
		err := h.telegramService.SendMessage(ctx, receiver.System.ExternalID,
			fmt.Sprintf("%s %s %s", child.Name.FirstName.String(), contextType, zone.Name.String()))
		if err != nil {
			h.observer.Logger.Trace().Err(err).Msg("invalid notification setting id")
			continue
		}
		notifySettingID, err := value_object.NewIDFromString(receiver.NotifySettingID.String())
		if err != nil {
			h.observer.Logger.Trace().Err(err).Msg("invalid notification setting id")
			continue
		}
		log, err := entity.NewNotificationLog(notifySettingID, notifyType, receiver.System.ID)
		if err != nil {
			h.observer.Logger.Trace().Err(err).Msg("failed to create notification log")
			continue
		}
		//fmt.Printf("notify log: %v\n", log.ID)
		if err = h.notificationLogStorage.AddNotificationLog(ctx, log); err != nil {
			if errors.Is(err, adapter.ErrStorage) {
				h.observer.Logger.Error().Err(err).Msg("database error while adding notification log")
				continue
			}
			h.observer.Logger.Error().Err(err).Msg("unexpected error while adding notification log")
			continue
		}
		//fmt.Printf("notify log: %v\n", log)
		notifyLogsIDs = append(notifyLogsIDs, log.ID)
	}
	return notifyLogsIDs, nil
}

func (h *AddUserLocationHandler) PushInsideNotifications(
	ctx context.Context,
	child *entity.User,
	zoneID value_object.ID,
	notifyType string,
) ([]value_object.ID, error) {
	zone, err := h.zoneStorage.GetZone(ctx, zoneID)
	if err != nil {
		if errors.Is(err, adapter.ErrZoneNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("zone not found")
			return []value_object.ID{}, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting zone")
			return []value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting zone")
		return []value_object.ID{}, err
	}
	settingIDs, err := h.notificationSettingStorage.GetDangerZoneNotificationSettings(ctx, child.ID, zoneID)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting zone notification settings")
			return []value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting zone notification settings")
		return []value_object.ID{}, err
	}
	if len(settingIDs) == 0 {
		h.observer.Logger.Info().Msg("no notification settings found for user location")
		return []value_object.ID{}, nil
	}
	//fmt.Printf("settingIDs: %v\n", settingIDs)

	receivers, err := h.notificationSettingStorage.GetIntegratedReceiversBySettingIDs(ctx, settingIDs)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting notification receivers")
			return []value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting notification receivers")
		return []value_object.ID{}, err
	}
	if len(receivers) == 0 {
		h.observer.Logger.Info().Msg("no integrated receivers found for user location")
		return []value_object.ID{}, nil
	}
	//fmt.Printf("receivers: %v\n", receivers)

	var settingIDsVO []value_object.ID
	for _, receiver := range receivers {
		settingID, err := value_object.NewIDFromString(receiver.NotifySettingID.String())
		if err != nil {
			h.observer.Logger.Trace().Err(err).Msg("invalid integrated user notification setting id")
			continue
		}
		settingIDsVO = append(settingIDsVO, settingID)
	}

	settings, err := h.notificationSettingStorage.GetNotificationSettingsByIDs(ctx, settingIDsVO)
	if err != nil {
		if errors.Is(err, adapter.ErrNotificationSettingNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("notification setting not found")
			return []value_object.ID{}, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting notification settings")
			return []value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting notification settings")
		return []value_object.ID{}, err
	}
	//fmt.Printf("settings: %v\n", settings)

	lastLogs, err := h.notificationLogStorage.GetLastLogsByNotificationIDs(ctx, settingIDsVO)
	if err != nil {
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting notification logs")
			return []value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting notification logs")
		return []value_object.ID{}, err
	}
	//fmt.Printf("lastLogs: %v\n", lastLogs)

	lastLogMap := make(map[string]*response.NotificationLogDB, len(lastLogs))
	for _, l := range lastLogs {
		lastLogMap[l.NotificationID] = l
	}
	//fmt.Printf("lastLogMap: %v\n", lastLogMap)
	loc, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(loc)

	settingsMap := make(map[string]*response.NotificationSettingWithFrequencyDB, len(settings))
	for _, s := range settings {
		settingsMap[s.ID] = s
	}

	var notificationLogsIDs []value_object.ID
	for _, receiver := range receivers {
		setting, ok := settingsMap[receiver.NotifySettingID.String()]
		if !ok {
			continue
		}
		last := lastLogMap[setting.ID]
		elapsed := time.Duration(0)
		if last != nil {
			elapsed = now.Sub(last.CreatedAt)
			//fmt.Printf("elapsed = %v\n", elapsed)
		}
		if setting.Frequency == nil {
			continue
		}
		if last == nil || elapsed >= *setting.Frequency {
			if receiver.System.ID != 1 {
				continue
			}
			var contextType string
			if notifyType == "inside" {
				contextType = "находится в"
			}
			err := h.telegramService.SendMessage(ctx, receiver.System.ExternalID,
				fmt.Sprintf("%s %s %s", child.Name.FirstName.String(), contextType, zone.Name.String()))
			if err != nil {
				h.observer.Logger.Trace().Err(err).Msg("invalid notification setting id")
				continue
			}
			fmt.Println("inside push")
			id, err := value_object.NewIDFromString(setting.ID)
			if err != nil {
				h.observer.Logger.Trace().Err(err).Msg("invalid notification log id")
				continue
			}
			log, err := entity.NewNotificationLog(id, notifyType, receiver.System.ID)
			if err != nil {
				h.observer.Logger.Trace().Err(err).Msg("failed to create notification log")
				continue
			}
			//fmt.Printf("notify log: %v\n", log.ID)
			if err := h.notificationLogStorage.AddNotificationLog(ctx, log); err != nil {
				if errors.Is(err, adapter.ErrStorage) {
					h.observer.Logger.Error().Err(err).Msg("database error while adding notification log")
					continue
				}
				h.observer.Logger.Error().Err(err).Msg("unexpected error while adding notification log")
				continue
			}
			notificationLogsIDs = append(notificationLogsIDs, log.ID)
		}
	}
	return notificationLogsIDs, nil
}

//settingIDs, err := h.notificationSettingStorage.GetDangerZoneNotificationSettings(ctx, user.ID, *ctxZone)
//if err != nil {
//	if errors.Is(err, adapter.ErrStorage) {
//		h.observer.Logger.Error().Err(err).Msg("database error while getting notification settings")
//		return value_object.ID{}, []value_object.ID{}, err
//	}
//	h.observer.Logger.Error().Err(err).Msg("unexpected error while getting notification settings")
//	return value_object.ID{}, []value_object.ID{}, err
//}
//if len(settingIDs) == 0 {
//	h.observer.Logger.Info().Msg("no notification settings found for user location")
//	return userLocation.ID, []value_object.ID{}, nil
//}
//fmt.Printf("settingIDs: %v\n", settingIDs)
//settings, err := h.notificationSettingStorage.GetNotificationSettingsByIDs(ctx, settingIDs)
//if err != nil {
//	if errors.Is(err, adapter.ErrNotificationSettingNotFound) {
//		h.observer.Logger.Trace().Err(err).Msg("notification setting not found")
//		return value_object.ID{}, []value_object.ID{}, err
//	}
//	if errors.Is(err, adapter.ErrStorage) {
//		h.observer.Logger.Error().Err(err).Msg("database error while getting notification settings")
//		return value_object.ID{}, []value_object.ID{}, err
//	}
//	h.observer.Logger.Error().Err(err).Msg("unexpected error while getting notification settings")
//	return value_object.ID{}, []value_object.ID{}, err
//}
//fmt.Printf("settings: %v\n", settings)
//lastLogs, err := h.notificationLogStorage.GetLastLogsByNotificationIDs(ctx, settingIDs)
//if err != nil {
//	if errors.Is(err, adapter.ErrStorage) {
//		h.observer.Logger.Error().Err(err).Msg("database error while getting notification logs")
//		return value_object.ID{}, []value_object.ID{}, err
//	}
//	h.observer.Logger.Error().Err(err).Msg("unexpected error while getting notification logs")
//	return value_object.ID{}, []value_object.ID{}, err
//}
//fmt.Printf("lastLogs: %v\n", lastLogs)
//lastLogMap := make(map[string]*response.NotificationLogDB, len(lastLogs))
//for _, l := range lastLogs {
//	lastLogMap[l.NotificationID] = l
//}
//fmt.Printf("lastLogMap: %v\n", lastLogMap)

//loc, _ := time.LoadLocation("Europe/Moscow")
//now := time.Now().In(loc)

//if notifyType == "entry" || notifyType == "exit" {
////var zoneForNotify *value_object.ID
////if notifyType == "entry" {
////	zoneForNotify = currZoneID
////} else {
////	zoneForNotify = prevZoneID
////}
//
//if zoneForNotify != nil {
//settingIDs, err := h.notificationSettingStorage.GetDangerZoneNotificationSettings(ctx, user.ID, zones[0].ID)
//if err != nil {
//if errors.Is(err, adapter.ErrStorage) {
//h.observer.Logger.Error().Err(err).Msg("database error while getting notification settings")
//return value_object.ID{}, []value_object.ID{}, err
//}
//h.observer.Logger.Error().Err(err).Msg("unexpected error while getting notification settings")
//return value_object.ID{}, []value_object.ID{}, err
//}
//if len(settingIDs) != 0 {
//var notifyLogsIDs []value_object.ID
//for _, settingID := range settingIDs {
//notifyLog, err := entity.NewNotificationLog(settingID, notifyType)
//if err != nil {
//h.observer.Logger.Trace().Err(err).Msg("failed to create notification log")
//return value_object.ID{}, []value_object.ID{}, fmt.Errorf("failed to create notification log: %w", err)
//}
//err = h.notificationLogStorage.AddNotificationLog(ctx, notifyLog)
//if err != nil {
//if errors.Is(err, adapter.ErrStorage) {
//h.observer.Logger.Error().Err(err).Msg("database error while adding notification log")
//return value_object.ID{}, []value_object.ID{}, err
//}
//h.observer.Logger.Error().Err(err).Msg("unexpected error while adding notification log")
//return value_object.ID{}, []value_object.ID{}, err
//}
//notifyLogsIDs = append(notifyLogsIDs, notifyLog.ID)
//}
//return userLocation.ID, notifyLogsIDs, nil
//}
//h.observer.Logger.Info().Msg("no notification settings found for user location")
//return userLocation.ID, []value_object.ID{}, err
//}
//}
//h.observer.Logger.Info().Msg("no push notifications")
//return userLocation.ID, []value_object.ID{}, err

//func (h *AddUserLocationHandler) checkZonesAndCreateNotificationLog(user *entity.User, locationTime *value_object.LocationTime) error {
//	memberships, err := h.membershipStorage.GetMembershipsByUserID(user.externalID)
//	if err != nil {
//		return err
//	}
//
//	for _, membership := range memberships {
//		familyZones, err := h.familyZonesStorage.GetZonesByFamilyID(membership.externalID)
//		if err != nil {
//			return err
//		}
//
//		familyID, err := value_object.NewIDFromString(membership.externalID)
//		if err != nil {
//			return err
//		}
//		allMemberships, err := h.membershipStorage.GetMembershipsByFamilyID(familyID)
//		if err != nil {
//			return err
//		}
//
//		for _, zoneID := range familyZones {
//			id, err := value_object.NewIDFromString(zoneID)
//			if err != nil {
//				return err
//			}
//			boundaries, err := h.boundariesStorage.GetBoundariesByZoneID(id)
//			if err != nil {
//				return err
//			}
//
//			zone := &entity.Zone{
//				externalID:         id,
//				Boundaries: boundaries,
//			}
//
//			previousLocationInZone := false
//			if len(user.Locations) > 1 {
//				previousLocation := user.Locations[len(user.Locations)-1].Location
//				previousLocationInZone = h.isLocationInZone(previousLocation, zone)
//			}
//			fmt.Println(fmt.Sprintf("Previous location in zone: %v", previousLocationInZone))
//
//			currentLocationInZone := h.isLocationInZone(locationTime.Location, zone)
//			fmt.Println(fmt.Sprintf("Current location in zone: %v", currentLocationInZone))
//			context := ""
//			if !previousLocationInZone && currentLocationInZone {
//				context = "entrance to the zone"
//			} else if previousLocationInZone && !currentLocationInZone {
//				context = "leaving the zone"
//			}
//			fmt.Println(fmt.Sprintf("Context: %v", context))
//			fmt.Println("\n\n")
//
//			if context != "" {
//				for _, m := range allMemberships {
//					if m.UserID != user.externalID.String() {
//						receiverID, err := value_object.NewIDFromString(m.UserID)
//						if err != nil {
//							return err
//						}
//
//						senderID := user.externalID
//						notificationSetting := h.notificationSettingStorage.GetNotificationSettingByRecieverSender(
//							receiverID, senderID)
//						fmt.Println(fmt.Sprintf("Notification setting: %v", notificationSetting.externalID))
//
//						if notificationSetting == nil {
//							fmt.Println("Notification setting is nil")
//							continue
//						}
//
//						notificationLog, err := aggregate.NewNotificationLog(&notificationSetting.externalID, time.Now(), context)
//						if err != nil {
//							return err
//						}
//						err = h.notificationLogStorage.AddNotificationLog(notificationLog)
//						if err != nil {
//							return err
//						}
//						fmt.Println(fmt.Sprintf("Notification log: %v", notificationLog.externalID))
//					}
//				}
//			}
//		}
//	}
//
//	return nil
//}
//
//func (h *AddUserLocationHandler) isLocationInZone(location *value_object.Location, zone *entity.Zone) bool {
//	boundaries := *zone.Boundaries
//	n := len(boundaries)
//	if n < 3 {
//		return false
//	}
//
//	lat := float64(*location.Latitude)
//	lon := float64(*location.Longitude)
//	inside := false
//
//	j := n - 1
//	for i := 0; i < n; i++ {
//		latI := float64(*boundaries[i].Latitude)
//		lonI := float64(*boundaries[i].Longitude)
//		latJ := float64(*boundaries[j].Latitude)
//		lonJ := float64(*boundaries[j].Longitude)
//
//		if (lonI > lon) != (lonJ > lon) &&
//			(lat < (latJ-latI)*(lon-lonI)/(lonJ-lonI)+latI) {
//			inside = !inside
//		}
//		j = i
//	}
//
//	return inside
//}
