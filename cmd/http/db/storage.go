package db

import (
	"awesomeProjectDDD/internal/adapter/db_storage"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/pkg/postgres"
	"github.com/gofiber/fiber/v2/log"
)

func CreateStorages(
	observer *observability.Observability,
	pg *postgres.Postgres,
) (
	*db_storage.UserRepositoryImpl,
	*db_storage.FamilyRepositoryImpl,
	*db_storage.ZoneRepositoryImpl,
	*db_storage.RoleRepositoryImpl,
	*db_storage.FamilyMembershipRepositoryImpl,
	*db_storage.InvitationToFamilyRepositoryImpl,
	//*db_storage.MessengerTypeRepositoryImpl,
	*db_storage.NotificationFrequencyRepositoryImpl,
	*db_storage.NotificationSettingRepositoryImpl,
	//*db_storage.NotificationTypeRepositoryImpl,
	//*db_storage.NotifySettingNotifyTypeRepositoryImpl,
	//*db_storage.NotifySettingMesTypeRepositoryImpl,
	*db_storage.UserLocationsRepositoryImpl,
	*db_storage.NotificationLogRepositoryImpl,
	*db_storage.IntegrationRepositoryImpl,
	error) {
	dbUser, err := db_storage.NewUserRepositoryImpl(observer, pg)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	dbFamily, err := db_storage.NewFamilyRepositoryImpl(pg, observer)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	dbZone, err := db_storage.NewZoneRepositoryImpl(pg, observer)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	dbRole, err := db_storage.NewRoleRepositoryImpl(pg, observer)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	dbFamilyMembership, err := db_storage.NewFamilyMembershipRepositoryImpl(pg, observer)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	dbInvitationToFamily, err := db_storage.NewInvitationToFamilyRepositoryImpl(pg, observer)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	//dbMessengerType, err := db_storage.NewMessengerTypeRepositoryImpl(pg, observer)
	//if err != nil {
	//	log.Fatalf("failed to connect to the database: %v", err)
	//}

	dbNotFrequency, err := db_storage.NewNotificationFrequencyRepositoryImpl(pg, observer)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	dbNotificationSetting, err := db_storage.NewNotificationSettingRepositoryImpl(pg, observer)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	//dbNotificationTypes, err := db_storage.NewNotificationTypeRepositoryImpl(pg, observer)
	//if err != nil {
	//	log.Fatalf("failed to connect to the database: %v", err)
	//}

	//dbNotSettingNotTypes, err := db_storage.NewNotifySettingNotifyTypeRepositoryImpl(pg, observer)
	//if err != nil {
	//	log.Fatalf("failed to connect to the database: %v", err)
	//}

	//dbNotSettingMesTypes, err := db_storage.NewNotifySettingMesTypeRepositoryImpl(pg, observer)
	//if err != nil {
	//	log.Fatalf("failed to connect to the database: %v", err)
	//}

	dbUserLocations, err := db_storage.NewUserLocationsRepositoryImpl(pg, observer)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	dbNotifyLog, err := db_storage.NewNotificationLogRepositoryImpl(pg, observer)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	dbIntegration, err := db_storage.NewIntegrationRepositoryImpl(pg, observer)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	return dbUser, dbFamily, dbZone, dbRole, dbFamilyMembership, dbInvitationToFamily,
		//dbMessengerType,
		dbNotFrequency,
		dbNotificationSetting,
		//dbNotificationTypes,
		//dbNotSettingMesTypes,
		dbUserLocations, dbNotifyLog, dbIntegration, nil
}
