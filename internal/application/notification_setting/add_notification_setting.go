package notification_setting

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/application"
	"awesomeProjectDDD/internal/domain"
	"awesomeProjectDDD/internal/domain/model/aggregate"
	"awesomeProjectDDD/internal/domain/model/value_object"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/port/db"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type AddNotificationSettingCommand struct {
	FrequencyID *uuid.UUID
	EventType   string
	ReceiverID  string
	SenderID    uuid.UUID
	ZoneID      *uuid.UUID
	Battery     *int
}

type AddNotificationSettingHandler struct {
	notificationSettingStorage db.NotificationSettingRepository
	frequencyStorage           db.NotificationFrequencyRepository
	//notificationTypeStorage    db.NotificationTypeRepository
	//notSettingNotTypesStorage  db.NotifySettingNotifyTypeRepository
	//messengerTypeStorage       db.MessengerTypeRepository
	//notSettingMesTypesStorage  db.NotifySettingMesTypeRepository
	userStorage       db.UserRepository
	familyStorage     db.FamilyRepository
	membershipStorage db.FamilyMembershipRepository
	roleStorage       db.RoleRepository
	zoneStorage       db.ZoneRepository
	observer          *observability.Observability
}

func NewAddNotificationSettingHandler(
	notificationSettingStorage db.NotificationSettingRepository,
	frequencyStorage db.NotificationFrequencyRepository,
	//notificationTypeStorage db.NotificationTypeRepository,
	//notSettingNotTypesStorage db.NotifySettingNotifyTypeRepository,
	//messengerTypeStorage db.MessengerTypeRepository,
	//notSettingMesTypesStorage db.NotifySettingMesTypeRepository,
	userStorage db.UserRepository,
	familyStorage db.FamilyRepository,
	membershipStorage db.FamilyMembershipRepository,
	roleStorage db.RoleRepository,
	zoneStorage db.ZoneRepository,
	observer *observability.Observability,
) *AddNotificationSettingHandler {
	return &AddNotificationSettingHandler{
		notificationSettingStorage: notificationSettingStorage,
		frequencyStorage:           frequencyStorage,
		//notificationTypeStorage:    notificationTypeStorage,
		//notSettingNotTypesStorage:  notSettingNotTypesStorage,
		//messengerTypeStorage:       messengerTypeStorage,
		//notSettingMesTypesStorage:  notSettingMesTypesStorage,
		userStorage:       userStorage,
		familyStorage:     familyStorage,
		membershipStorage: membershipStorage,
		roleStorage:       roleStorage,
		zoneStorage:       zoneStorage,
		observer:          observer,
	}
}

func (h *AddNotificationSettingHandler) Handle(ctx context.Context,
	cmd AddNotificationSettingCommand) (value_object.ID, error) {
	extID, err := value_object.NewIDFromString(cmd.ReceiverID)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid user external id")
		return value_object.ID{}, err
	}
	receiver, err := h.userStorage.GetUserByExternalID(ctx, extID)
	if err != nil {
		if errors.Is(err, adapter.ErrUserNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("notify receiver not found")
			return value_object.ID{}, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting notify receiver")
			return value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting notify receiver")
		return value_object.ID{}, err
	}
	if receiver.ID.ToRaw() == cmd.SenderID {
		h.observer.Logger.Trace().Msg("receiver is sender")
		return value_object.ID{}, application.ErrReceiverIsSender
	}

	senderID, err := value_object.NewIDFromString(cmd.SenderID.String())
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("invalid sender id")
		return value_object.ID{}, err
	}
	sender, err := h.userStorage.GetUser(ctx, senderID)
	if err != nil {
		if errors.Is(err, adapter.ErrUserNotFound) {
			h.observer.Logger.Trace().Err(err).Msg("notify sender not found")
			return value_object.ID{}, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while getting notify sender")
			return value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while getting notify sender")
		return value_object.ID{}, err
	}
	var notifySettingID value_object.ID
	switch cmd.EventType {
	case "zone":
		fmt.Printf("case zone")
		if cmd.ZoneID == nil {
			h.observer.Logger.Error().Err(err).Msg("zone id is required for zone event")
			return value_object.ID{}, domain.ErrInvalidZone
		}
		zoneID, err := value_object.NewIDFromString(cmd.ZoneID.String())
		if err != nil {
			h.observer.Logger.Trace().Err(err).Msg("invalid zone id")
			return value_object.ID{}, err
		}
		zone, err := h.zoneStorage.GetZone(ctx, zoneID)
		if err != nil {
			if errors.Is(err, adapter.ErrZoneNotFound) {
				h.observer.Logger.Trace().Err(err).Msg("zone not found")
				return value_object.ID{}, err
			}
			if errors.Is(err, adapter.ErrStorage) {
				h.observer.Logger.Error().Err(err).Msg("database error while getting zone")
				return value_object.ID{}, err
			}
			h.observer.Logger.Error().Err(err).Msg("unexpected error while getting zone")
			return value_object.ID{}, err
		}
		belong, err := h.membershipStorage.CheckUsersBelongToFamilyByZone(ctx, zone.ID, receiver.ID, sender.ID)
		if err != nil {
			if errors.Is(err, adapter.ErrMembershipNotFound) {
				h.observer.Logger.Trace().Err(err).Msg("membership not found")
				return value_object.ID{}, err
			}
			if errors.Is(err, adapter.ErrNotFamilyMembers) {
				h.observer.Logger.Trace().Err(err).Msg("not family members")
				return value_object.ID{}, err
			}
			if errors.Is(err, adapter.ErrStorage) {
				h.observer.Logger.Error().Err(err).Msg("database error while checking belong to one family")
				return value_object.ID{}, err
			}
			h.observer.Logger.Error().Err(err).Msg("unexpected error while checking belong to one family")
			return value_object.ID{}, err
		}
		if belong {
			if zone.Safety.Bool() == false {
				if cmd.FrequencyID == nil {
					h.observer.Logger.Error().Err(err).Msg("frequency id is required for dangerous zone")
					return value_object.ID{}, domain.ErrInvalidFrequency
				}
				frequencyID, err := value_object.NewIDFromString(cmd.FrequencyID.String())
				if err != nil {
					h.observer.Logger.Trace().Err(err).Msg("invalid frequency id")
					return value_object.ID{}, err
				}
				frequency, err := h.frequencyStorage.GetFrequency(ctx, frequencyID)
				if err != nil {
					if errors.Is(err, adapter.ErrFrequencyNotFound) {
						h.observer.Logger.Trace().Err(err).Msg("frequency not found")
						return value_object.ID{}, err
					}
					if errors.Is(err, adapter.ErrStorage) {
						h.observer.Logger.Error().Err(err).Msg("database error while getting notify frequency")
						return value_object.ID{}, err
					}
					h.observer.Logger.Error().Err(err).Msg("unexpected error while getting notify frequency")
					return value_object.ID{}, err
				}
				notifySetting, err := aggregate.NewNotificationSetting(frequency, cmd.EventType, receiver, sender, zone, nil)
				if err != nil {
					h.observer.Logger.Trace().Err(err).Msg("failed to create notify setting")
					return value_object.ID{}, fmt.Errorf("failed to create notify setting: %w", err)
				}
				notifySettingID, err = h.InsertZoneSetting(ctx, notifySetting, zone.ID)
				if err != nil {
					h.observer.Logger.Trace().Err(err).Msg("failed to create notify setting")
					return value_object.ID{}, fmt.Errorf("failed to create notify setting: %w", err)
				}
			} else {
				notifySetting, err := aggregate.NewNotificationSetting(nil, cmd.EventType, receiver, sender, zone, nil)
				if err != nil {
					h.observer.Logger.Trace().Err(err).Msg("failed to create notify setting")
					return value_object.ID{}, fmt.Errorf("failed to create notify setting: %w", err)
				}
				notifySettingID, err = h.InsertZoneSetting(ctx, notifySetting, zone.ID)
				if err != nil {
					h.observer.Logger.Trace().Err(err).Msg("failed to create notify setting")
					return value_object.ID{}, fmt.Errorf("failed to create notify setting: %w", err)
				}
			}
		} else {
			h.observer.Logger.Trace().Err(err).Msg("not family members")
			return value_object.ID{}, adapter.ErrNotFamilyMembers
		}

	case "battery":
		fmt.Printf("case battery")
		if cmd.Battery == nil {
			h.observer.Logger.Error().Err(err).Msg("min battery is required for battery event")
			return value_object.ID{}, domain.ErrInvalidBattery
		}
		if cmd.FrequencyID == nil {
			h.observer.Logger.Error().Err(err).Msg("frequency id is required for battery event")
			return value_object.ID{}, domain.ErrInvalidFrequency
		}
		battery, err := value_object.NewBatteryThreshold(*cmd.Battery)
		if err != nil {
			h.observer.Logger.Trace().Err(err).Msg("failed to add battery threshold")
			return value_object.ID{}, fmt.Errorf("failed to add battery threshold: %w", err)
		}
		belong, err := h.membershipStorage.CheckUsersShareCommonActiveFamily(ctx, sender.ID, receiver.ID)
		if err != nil {
			if errors.Is(err, adapter.ErrMembershipNotFound) {
				h.observer.Logger.Trace().Err(err).Msg("membership not found")
				return value_object.ID{}, err
			}
			if errors.Is(err, adapter.ErrNotFamilyMembers) {
				h.observer.Logger.Trace().Err(err).Msg("not family members")
				return value_object.ID{}, err
			}
			if errors.Is(err, adapter.ErrStorage) {
				h.observer.Logger.Error().Err(err).Msg("database error while checking belong to one family")
				return value_object.ID{}, err
			}
			h.observer.Logger.Error().Err(err).Msg("unexpected error while checking belong to one family")
			return value_object.ID{}, err
		}
		if belong {
			frequencyID, err := value_object.NewIDFromString(cmd.FrequencyID.String())
			if err != nil {
				h.observer.Logger.Trace().Err(err).Msg("invalid frequency id")
				return value_object.ID{}, err
			}
			frequency, err := h.frequencyStorage.GetFrequency(ctx, frequencyID)
			if err != nil {
				if errors.Is(err, adapter.ErrFrequencyNotFound) {
					h.observer.Logger.Trace().Err(err).Msg("frequency not found")
					return value_object.ID{}, err
				}
				if errors.Is(err, adapter.ErrStorage) {
					h.observer.Logger.Error().Err(err).Msg("database error while getting notify frequency")
					return value_object.ID{}, err
				}
				h.observer.Logger.Error().Err(err).Msg("unexpected error while getting notify frequency")
				return value_object.ID{}, err
			}
			notifySetting, err := aggregate.NewNotificationSetting(frequency, cmd.EventType, receiver, sender, nil, &battery)
			if err != nil {
				h.observer.Logger.Trace().Err(err).Msg("failed to create notify setting")
				return value_object.ID{}, fmt.Errorf("failed to create notify setting: %w", err)
			}
			notifySettingID, err = h.InsertBatterySetting(ctx, notifySetting, *cmd.Battery)
			if err != nil {
				h.observer.Logger.Trace().Err(err).Msg("failed to create notify setting")
				return value_object.ID{}, fmt.Errorf("failed to create notify setting: %w", err)
			}
		} else {
			h.observer.Logger.Trace().Err(err).Msg("not family members")
			return value_object.ID{}, adapter.ErrNotFamilyMembers
		}
	default:
		h.observer.Logger.Error().Err(err).Msg("unknown event type")
		return value_object.ID{}, domain.ErrInvalidEventType
	}
	return notifySettingID, nil
}

func (h *AddNotificationSettingHandler) InsertZoneSetting(
	ctx context.Context,
	notifySetting *aggregate.NotificationSetting,
	zoneID value_object.ID,
) (value_object.ID, error) {
	exists, err := h.notificationSettingStorage.NotificationSettingExists(ctx, notifySetting)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidZoneOrBattery) {
			h.observer.Logger.Error().Err(err).Msg("zone or battery are required")
			return value_object.ID{}, err
		}
		if errors.Is(err, adapter.ErrNotificationSettingExists) {
			h.observer.Logger.Error().Err(err).Msg("notify setting already exists")
			return value_object.ID{}, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while checking notify setting exists")
			return value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while checking notify setting exists")
		return value_object.ID{}, err
	}
	if exists {
		h.observer.Logger.Trace().Msg("notify setting already exists in system")
		return value_object.ID{}, adapter.ErrNotificationSettingExists
	}

	notifySettingID := notifySetting.ID
	err = h.notificationSettingStorage.AddNotificationSetting(ctx, notifySetting)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to create notify setting")
		return value_object.ID{}, fmt.Errorf("failed to create notify setting: %w", err)
	}
	zoneNotificationSettingID := value_object.NewID()
	err = h.notificationSettingStorage.AddZoneNotificationSetting(ctx, zoneNotificationSettingID, notifySetting.ID, zoneID)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to create zone notify setting")
		return value_object.ID{}, fmt.Errorf("failed to create zone notify setting: %w", err)
	}
	return notifySettingID, nil
}

func (h *AddNotificationSettingHandler) InsertBatterySetting(
	ctx context.Context,
	notifySetting *aggregate.NotificationSetting,
	threshold int,
) (value_object.ID, error) {
	exists, err := h.notificationSettingStorage.NotificationSettingExists(ctx, notifySetting)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidZoneOrBattery) {
			h.observer.Logger.Error().Err(err).Msg("zone or battery are required")
			return value_object.ID{}, err
		}
		if errors.Is(err, adapter.ErrNotificationSettingExists) {
			h.observer.Logger.Error().Err(err).Msg("notify setting already exists")
			return value_object.ID{}, err
		}
		if errors.Is(err, adapter.ErrStorage) {
			h.observer.Logger.Error().Err(err).Msg("database error while checking notify setting exists")
			return value_object.ID{}, err
		}
		h.observer.Logger.Error().Err(err).Msg("unexpected error while checking notify setting exists")
		return value_object.ID{}, err
	}
	if exists {
		h.observer.Logger.Trace().Msg("notify setting already exists in system")
		return value_object.ID{}, adapter.ErrNotificationSettingExists
	}

	notifySettingID := notifySetting.ID
	err = h.notificationSettingStorage.AddNotificationSetting(ctx, notifySetting)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to create notify setting")
		return value_object.ID{}, fmt.Errorf("failed to create notify setting: %w", err)
	}
	zoneNotificationSettingID := value_object.NewID()
	err = h.notificationSettingStorage.AddBatteryNotificationSetting(ctx, zoneNotificationSettingID, notifySetting.ID, threshold)
	if err != nil {
		h.observer.Logger.Trace().Err(err).Msg("failed to create battery notify setting")
		return value_object.ID{}, fmt.Errorf("failed to create battery notify setting: %w", err)
	}
	return notifySettingID, nil
}

//var NotifyTypesIDs []value_object.ID
//for _, id := range cmd.NotifyTypesIDs {
//	notTypeID, err := value_object.NewIDFromString(id)
//	if err != nil {
//		h.observer.Logger.Trace().Err(err).Msg("invalid notify type id")
//		return "", fmt.Errorf("invalid notify type id: %w", err)
//	}
//	NotifyTypesIDs = append(NotifyTypesIDs, notTypeID)
//}
//_, err = h.notificationTypeStorage.GetNotificationTypesByIDs(ctx, cmd.NotifyTypesIDs)
//if err != nil {
//	h.observer.Logger.Error().Err(err).Msg("failed to get notification types")
//	return "", fmt.Errorf("failed to get notification types: %w", err)
//}

//var MesTypesIDs []value_object.ID
//for _, id := range cmd.MesTypesIDs {
//	mesTypeID, err := value_object.NewIDFromString(id)
//	if err != nil {
//		h.observer.Logger.Trace().Err(err).Msg("invalid messenger type id")
//		return "", fmt.Errorf("invalid messenger type id: %w", err)
//	}
//	MesTypesIDs = append(MesTypesIDs, mesTypeID)
//}
//mesTypes, err := h.messengerTypeStorage.GetMessengerTypesByIDs(ctx, cmd.MesTypesIDs)
//if err != nil {
//	h.observer.Logger.Error().Err(err).Msg("failed to get messenger types")
//	return "", fmt.Errorf("failed to get messenger types: %w", err)
//}

//receiverID, err := value_object.NewIDFromString(cmd.ReceiverID)
//if err != nil {
//	h.observer.Logger.Trace().Err(err).Msg("invalid receiver id")
//	return "", fmt.Errorf("invalid receiver id: %w", err)
//}
//senderID, err := value_object.NewIDFromString(cmd.SenderID)
//if err != nil {
//	h.observer.Logger.Trace().Err(err).Msg("invalid sender id")
//	return "", fmt.Errorf("invalid sender id: %w", err)
//}
//_, err = h.userStorage.GetUsersByIDs(ctx, []string{receiverID.String(), senderID.String()})
//if err != nil {
//	h.observer.Logger.Error().Err(err).Msg("failed to get sender and receiver")
//	return "", fmt.Errorf("failed to get sender and receiver: %w", err)
//}
//
//zoneID, err := value_object.NewIDFromString(cmd.ZoneID)
//if err != nil {
//	h.observer.Logger.Trace().Err(err).Msg("invalid zone id")
//	return "", fmt.Errorf("invalid zone id: %w", err)
//}
//_, err = h.zoneStorage.GetZone(ctx, zoneID)
//if err != nil {
//	h.observer.Logger.Error().Err(err).Msg("failed to get zone")
//	return "", fmt.Errorf("failed to get zone: %w", err)
//}
//
//familyID, err := value_object.NewIDFromString(cmd.FamilyID)
//if err != nil {
//	h.observer.Logger.Trace().Err(err).Msg("invalid family id")
//	return "", fmt.Errorf("invalid family id: %w", err)
//}
//_, err = h.familyStorage.GetFamily(ctx, familyID)
//if err != nil {
//	h.observer.Logger.Error().Err(err).Msg("failed to get family")
//	return "", fmt.Errorf("failed to get family: %w", err)
//}
//
//memberships, err := h.membershipStorage.GetMembershipsByFamilyID(ctx, familyID,
//	&[]string{receiverID.String(), senderID.String()})
//if err != nil {
//	h.observer.Logger.Error().Err(err).Msg("failed to get memberships")
//	return "", fmt.Errorf("failed to get memberships: %w", err)
//}
//if len(memberships) != 2 {
//	h.observer.Logger.Trace().Msg("receiver and sender must be in the family")
//	return "", fmt.Errorf("receiver and sender must be in the family")
//}
//receiverMembershipID, err := value_object.NewIDFromString(memberships[0].externalID)
//if err != nil {
//	h.observer.Logger.Trace().Err(err).Msg("invalid receiver membership id")
//	return "", fmt.Errorf("invalid receiver membership id: %w", err)
//}
//receiverRoleID, err := strconv.Atoi(memberships[0].RoleID)
//if err != nil {
//	h.observer.Logger.Trace().Err(err).Msg("invalid receiver role id")
//	return "", fmt.Errorf("invalid receiver role id: %w", err)
//}
//receiverRole, err := h.roleStorage.GetRole(ctx, receiverRoleID)
//if err != nil {
//	h.observer.Logger.Error().Err(err).Msg("failed to get receiver role")
//	return "", fmt.Errorf("failed to get receiver role: %w", err)
//}
//receiverMembership := aggregate.FamilyMembership{
//	externalID:        receiverMembershipID,
//	User:      users[1],
//	Role:      receiverRole,
//	Family:    family,
//	CreatedAt: memberships[0].CreatedAt,
//}

//senderMembershipID, err := value_object.NewIDFromString(memberships[1].externalID)
//if err != nil {
//	h.observer.Logger.Trace().Err(err).Msg("invalid sender membership id")
//	return "", fmt.Errorf("invalid sender membership id: %w", err)
//}
//senderRoleID, err := strconv.Atoi(memberships[1].RoleID)
//if err != nil {
//	h.observer.Logger.Trace().Err(err).Msg("invalid sender role id")
//	return "", fmt.Errorf("invalid sender role id: %w", err)
//}
//senderRole, err := h.roleStorage.GetRole(ctx, senderRoleID)
//if err != nil {
//	h.observer.Logger.Error().Err(err).Msg("failed to get sender role")
//	return "", fmt.Errorf("failed to get sender role: %w", err)
//}
//senderMembership := aggregate.FamilyMembership{
//	externalID:        senderMembershipID,
//	User:      users[0],
//	Role:      senderRole,
//	Family:    family,
//	CreatedAt: memberships[1].CreatedAt,
//}

//notificationSetting, err := aggregate.NewNotificationSetting(
//	frequency,
//	notifyTypes,
//	//mesTypes,
//	users[1],
//	users[0],
//	zone,
//	family)
//if err != nil {
//	h.observer.Logger.Trace().Err(err).Msg("failed to create notification setting aggregate")
//	return "", fmt.Errorf("failed to create notification setting aggregate: %w", err)
//}
//err = h.notificationSettingStorage.CreateNotificationSetting(ctx, notificationSetting)
//if err != nil {
//	h.observer.Logger.Error().Err(err).Msg("failed to save notification setting")
//	return "", fmt.Errorf("failed to save notification setting: %w", err)
//}

//var notifySettingNotTypes []*entity.NotifySettingNotifyType
//for _, notifyType := range notifyTypes {
//	notSettingNotType, err := entity.NewNotifySettingNotifyType(notificationSetting.ID, notifyType.ID)
//	if err != nil {
//		h.observer.Logger.Trace().Err(err).Msg("failed to create notify setting and notify type relation")
//		return "", fmt.Errorf("failed to create notify setting and notify type relation: %w", err)
//	}
//	notifySettingNotTypes = append(notifySettingNotTypes, notSettingNotType)
//}
//
//err = h.notSettingNotTypesStorage.AddNotifySettingNotTypes(ctx, notifySettingNotTypes)
//if err != nil {
//	h.observer.Logger.Error().Err(err).Msg("failed to add notify setting and notify type relation to db")
//	return "", fmt.Errorf("failed to add notify setting and notify type relation to db: %w", err)
//}
//
//var notifySettingMesTypes []*entity.NotifySettingMesType
//for _, mesType := range mesTypes {
//	notSettingMesType, err := entity.NewNotifySettingMesType(notificationSetting.ID, mesType.ID)
//	if err != nil {
//		h.observer.Logger.Trace().Err(err).Msg("failed to create notify setting and messenger type relation")
//		return "", fmt.Errorf("failed to create notify setting and messenger type relation: %w", err)
//	}
//	notifySettingMesTypes = append(notifySettingMesTypes, notSettingMesType)
//}
//err = h.notSettingMesTypesStorage.AddNotifySettingMesTypes(ctx, notifySettingMesTypes)
//if err != nil {
//	h.observer.Logger.Error().Err(err).Msg("failed to add notify setting and messenger type relation to db")
//	return "", fmt.Errorf("failed to add notify setting and messenger type relation to db: %w", err)
//}
//
//return notificationSetting.ID.String(), nil
//return "", nil
//}
