package main

import (
	"awesomeProjectDDD/cmd/http/db"
	"awesomeProjectDDD/config"
	_ "awesomeProjectDDD/docs"
	"awesomeProjectDDD/internal/adapter/hydra"
	"awesomeProjectDDD/internal/adapter/telegram"
	"awesomeProjectDDD/internal/application/family"
	"awesomeProjectDDD/internal/application/integration"
	invitationToFamilyApp "awesomeProjectDDD/internal/application/invitation_to_family"
	membershipApp "awesomeProjectDDD/internal/application/membership"
	"awesomeProjectDDD/internal/application/notification_frequency"
	notificationSettingApp "awesomeProjectDDD/internal/application/notification_setting"
	"awesomeProjectDDD/internal/application/role"
	"awesomeProjectDDD/internal/application/user"
	"awesomeProjectDDD/internal/application/user_location"
	"awesomeProjectDDD/internal/application/zone"
	"awesomeProjectDDD/internal/observability"
	"awesomeProjectDDD/internal/presentation/http"
	familyPres "awesomeProjectDDD/internal/presentation/http/family"
	integration2 "awesomeProjectDDD/internal/presentation/http/integration"
	invitationToFamilyPres "awesomeProjectDDD/internal/presentation/http/invitation_to_family"
	membershipPres "awesomeProjectDDD/internal/presentation/http/membership"
	notificationFrequencyPres "awesomeProjectDDD/internal/presentation/http/notification_frequency"
	notificationSettingPres "awesomeProjectDDD/internal/presentation/http/notification_setting"
	rolePres "awesomeProjectDDD/internal/presentation/http/role"
	userPres "awesomeProjectDDD/internal/presentation/http/user"
	userLocationPres "awesomeProjectDDD/internal/presentation/http/user_location"
	zonePres "awesomeProjectDDD/internal/presentation/http/zone"
	"awesomeProjectDDD/pkg/docs/rapidoc"
	"awesomeProjectDDD/pkg/docs/swagger"
	httpResty "awesomeProjectDDD/pkg/http"
	"awesomeProjectDDD/pkg/logging"
	"awesomeProjectDDD/pkg/postgres"
	"fmt"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"html/template"
	"time"
)

// @title            Diplom Documentation
// @version          1.0
// @description        This is a sample swagger for diplom
// @host            localhost:3000
// @BasePath          /
// @securityDefinitions.apikey	APIKeyAuth
// @in							header
// @name						Authorization
// @description					Access token from SSO

func main() {
	cfg := config.New()

	logger := logging.NewZeroLogger(zerolog.TraceLevel)
	observer := observability.New(logger)

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	app.Use(compress.New(compress.Config{Level: compress.LevelBestSpeed}))

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: observer.GetLogger(),
		Fields: []string{
			fiberzerolog.FieldIP,
			fiberzerolog.FieldUserAgent,
			fiberzerolog.FieldLatency,
			fiberzerolog.FieldStatus,
			fiberzerolog.FieldMethod,
			fiberzerolog.FieldURL,
			fiberzerolog.FieldError,
			fiberzerolog.FieldBytesReceived,
			fiberzerolog.FieldBytesSent,
		},
	}))
	app.Get("/swagger/*", swagger.New(swagger.Config{
		Title:                  "Diplom Documentation",
		InstanceName:           "v1",
		URL:                    "v1/swagger/openapi.json",
		DeepLinking:            false,
		DisplayRequestDuration: true,
		Presets: []template.JS{
			"SwaggerUIBundle.presets.apis",
			"SwaggerUIStandalonePreset",
		},
		Layout: "BaseLayout",
	}))
	app.Get("/rapidoc/*", rapidoc.New(rapidoc.Config{
		Title:      "Diplom Documentation",
		HeaderText: "Diplom Documentation",
		SpecURL:    "swagger/openapi.json",
		LogoURL:    "https://rapidocweb.com/images/logo.png",
		Theme:      rapidoc.Theme_Light,
	}))

	pg, err := postgres.New(
		cfg.Postgres.Dsn(),
		postgres.MaxPoolSize(10),
		postgres.ConnAttempts(20),
		postgres.ConnTimeout(5*time.Second),
	)
	if err != nil {
		panic(err)
	}

	httpClient := httpResty.NewClient(cfg.API.HydraAdminURL)
	hydraService := hydra.NewService(httpClient)
	httpClientTg := httpResty.NewClient(
		fmt.Sprintf(
			"%s%s",
			"https://api.telegram.org/bot",
			cfg.API.TelegramToken,
		),
	)
	telegramService := telegram.NewService(httpClientTg)

	dbUser, dbFamily, dbZone, dbRole, dbMembership, dbInvitationToFamily,
		//dbMessengerType,
		dbNotificationFrequency, dbNotificationSetting,
		//dbNotificationType, dbNotSettingNotType, dbNotSettingMesType,
		dbUserLocations, dbNotificationLog, dbIntegration, err := db.CreateStorages(observer, pg)

	addUserCommand := user.NewAddUserHandler(dbUser, observer)
	addUserHandler := userPres.NewAddUserHandler(addUserCommand, observer)
	//getUserHandler := userPres.NewGetUserHandler(getUserQuery, observer)
	getUsersCommand := user.NewGetAllUsersHandler(dbUser, observer)
	getUsersHandler := userPres.NewGetAllUsersHandler(getUsersCommand, observer)

	getUserQuery := user.NewGetUserHandler(observer, dbUser)
	getMeHandler := userPres.NewGetMeHandler(getUserQuery, observer)
	//changeUserCommand := user.NewChangeUserHandler(dbUser, observer)
	//changeUserHandler := userPres.NewChangeUserHandler(changeUserCommand, observer)
	changeUserTrackingCommand := user.NewChangeUserTrackingHandler(dbUser, observer)
	changeUserTrackingHandler := userPres.NewChangeUserTrackingHandler(changeUserTrackingCommand, observer)
	//deleteUserCommand := user.NewDeleteUserHandler(dbUser, dbMembership, observer)
	//deleteUserHandler := userPres.NewDeleteUserHandler(deleteUserCommand, observer)
	restoreUserCommand := user.NewRestoreUserHandler(dbUser)
	restoreUserHandler := userPres.NewRestoreUserHandler(restoreUserCommand)

	addFamilyCommand := family.NewAddFamilyHandler(dbFamily, dbUser, dbMembership, dbRole, observer)
	addFamilyHandler := familyPres.NewAddFamilyHandler(addFamilyCommand, observer)
	getFamilyCommand := family.NewGetFamilyHandler(dbFamily)
	getFamilyHandler := familyPres.NewGetFamilyHandler(getFamilyCommand)
	getUserFamiliesCommand := family.NewGetUserFamiliesHandler(dbFamily, dbUser)
	getUserFamiliesHandler := familyPres.NewGetUserFamiliesHandler(getUserFamiliesCommand)
	changeFamilyCommand := family.NewChangeFamilyHandler(dbFamily, observer)
	changeFamilyHandler := familyPres.NewChangeFamilyHandler(changeFamilyCommand, observer)
	deleteFamilyCommand := family.NewDeleteFamilyHandler(dbFamily)
	deleteFamilyHandler := familyPres.NewDeleteFamilyHandler(deleteFamilyCommand)

	//addZoneCommand := zone.NewAddZoneHandler(dbZone, dbFamily, observer)
	//addZoneHandler := zonePres.NewAddZoneHandler(addZoneCommand, observer)
	getZoneCommand := zone.NewGetZoneHandler(dbZone)
	getZoneHandler := zonePres.NewGetZoneHandler(getZoneCommand)
	getAllZonesCommand := zone.NewGetAllZonesHandler(dbZone)
	getAllZonesHandler := zonePres.NewGetAllZonesHandler(getAllZonesCommand)
	//getZonesByUserIDCommand := zone.NewGetZonesByUserIDHandler(dbZone, dbUser, observer)
	//getZonesByUserIDHandler := zonePres.NewGetZonesByUserIDHandler(getZonesByUserIDCommand, observer)
	changeZoneCommand := zone.NewChangeZoneHandler(dbZone, dbUser, observer)
	changeZoneHandler := zonePres.NewChangeZoneHandler(changeZoneCommand, observer)
	deleteZoneCommand := zone.NewDeleteZoneHandler(dbZone, dbUser, observer)
	deleteZoneHandler := zonePres.NewDeleteZoneHandler(deleteZoneCommand, observer)

	addFamilyZoneCommand := family.NewAddFamilyZoneHandler(
		dbFamily, dbZone, observer)
	addFamilyZoneHandler := familyPres.NewAddFamilyZoneHandler(addFamilyZoneCommand, observer)
	getFamilyZonesCommand := family.NewGetFamilyZonesHandler(dbFamily, observer)
	getFamilyZonesHandler := familyPres.NewGetFamilyZonesHandler(getFamilyZonesCommand, observer)
	//deleteFamilyZoneCommand := family.NewDeleteFamilyZoneHandler(dbFamily, dbZone, observer)
	//deleteFamilyZoneHandler := familyPres.NewDeleteFamilyZoneHandler(deleteFamilyZoneCommand, observer)
	//deleteFamilyZonesCommand := family.NewDeleteFamilyZonesHandler(dbFamily)
	//deleteFamilyZonesHandler := familyPres.NewDeleteFamilyZonesHandler(deleteFamilyZonesCommand)

	//addRoleCommand := role.NewAddRoleHandler(dbRole)
	//addRoleHandler := rolePres.NewAddRoleHandler(addRoleCommand)
	getRoleCommand := role.NewGetRoleHandler(dbRole, hydraService)
	getRoleHandler := rolePres.NewGetRoleHandler(getRoleCommand)
	getAllRolesCommand := role.NewGetAllRolesHandler(dbRole)
	getAllRolesHandler := rolePres.NewGetAllRolesHandler(getAllRolesCommand)

	getUserMembershipsCommand := membershipApp.NewGetUserMembershipsHandler(dbMembership, dbUser, dbFamily, observer)
	getUserMembershipsHandler := membershipPres.NewGetUserMembershipsHandler(getUserMembershipsCommand, observer)
	getMembershipsByFamilyCommand := membershipApp.NewGetMembershipsByFamilyHandler(dbMembership, dbUser, dbRole, dbFamily)
	getMembershipsByFamilyHandler := membershipPres.NewGetMembershipsByFamilyHandler(getMembershipsByFamilyCommand)
	getInvitationsToFamilyByReceiverIdCommand := invitationToFamilyApp.
		NewGetInvitationsToFamilyByReceiverIdHandler(dbInvitationToFamily, dbUser, dbFamily, dbRole, observer)
	getInvitationsToFamilyByReceiverIdHandler := invitationToFamilyPres.NewGetInvitationsToFamilyByReceiverIdHandler(
		getInvitationsToFamilyByReceiverIdCommand, observer)
	addUserMembershipCommand := membershipApp.NewAddUserMembershipHandler(
		dbMembership, dbUser, dbRole, dbFamily, dbInvitationToFamily, observer)
	addUserMembershipHandler := membershipPres.NewAddUserMembershipHandler(addUserMembershipCommand, observer)
	addInvitationToFamilyCommand := invitationToFamilyApp.NewAddInvitationToFamilyHandler(
		dbInvitationToFamily, dbUser, dbRole, dbFamily, observer)
	addInvitationToFamilyHandler := invitationToFamilyPres.NewAddInvitationToFamilyHandler(addInvitationToFamilyCommand, observer)
	addMembershipByInvitationCommand := membershipApp.NewAddMembershipByInvitationHandler(
		dbInvitationToFamily, dbMembership, dbUser, dbRole, dbFamily)
	addMembershipByInvitationHandler := membershipPres.NewAddMembershipByInvitationHandler(addMembershipByInvitationCommand)
	changeMembershipCommand := membershipApp.NewChangeFamilyMembershipHandler(dbFamily, dbMembership, observer)
	changeFamilyMembershipHandler := membershipPres.NewChangeFamilyMembershipHandler(changeMembershipCommand, observer)
	deleteUserMembershipCommand := membershipApp.NewDeleteUserMembershipHandler(dbMembership, dbUser, observer)
	deleteUserMembershipHandler := membershipPres.NewDeleteUserMembershipHandler(deleteUserMembershipCommand, observer)

	addUserLocationsCommand := user_location.NewAddUserLocationHandler(dbUserLocations, dbUser, dbNotificationLog,
		dbFamily, dbMembership, dbZone, dbNotificationSetting, dbNotificationFrequency, observer, telegramService)
	addUserLocationsHandler := userLocationPres.NewAddUserLocationHandler(addUserLocationsCommand, observer, telegramService)
	getUserLocationsCommand := user_location.NewGetUserLocationsHandler(dbUser, dbNotificationSetting, dbUserLocations, observer)
	getUserLocationsHandler := userLocationPres.NewGetUserLocationsHandler(getUserLocationsCommand, observer)

	getAllNotificationFrequenciesCommand := notification_frequency.NewGetAllNotificationFrequenciesHandler(
		dbNotificationFrequency, observer)
	getAllNotificationFrequenciesHandler := notificationFrequencyPres.NewGetAllNotificationFrequenciesHandler(
		getAllNotificationFrequenciesCommand, observer)
	//getAllMessengerTypesCommand := messenger_type.NewGetAllMessengerTypesHandler(dbMessengerType, observer)
	//getAllMessengerTypesHandler := messengertype2.NewGetAllMessengerTypesHandler(
	//	getAllMessengerTypesCommand, observer)

	//getAllNotificationTypesCommand := notification_type.NewGetAllNotificationTypesHandler(dbNotificationType, observer)
	//getAllNotificationTypesHandler := notificationtype2.NewGetAllNotificationTypesHandler(
	//	getAllNotificationTypesCommand, observer)

	getAvailableNotificationSendersCommand := notificationSettingApp.NewGetAvailableNotificationSendersHandler(dbUser, dbFamily,
		dbMembership, observer)
	getAvailableNotificationSendersHandler := notificationSettingPres.NewGetAvailableNotificationSendersHandler(
		getAvailableNotificationSendersCommand, observer)

	addNotificationSettingCommand := notificationSettingApp.NewAddNotificationSettingHandler(
		dbNotificationSetting, dbNotificationFrequency, dbUser, dbFamily, dbMembership, dbRole, dbZone, observer)
	addNotificationSettingHandler := notificationSettingPres.NewAddNotificationSettingHandler(
		addNotificationSettingCommand, observer)
	getNotificationSettingsCommand := notificationSettingApp.NewGetNotificationSettingsHandler(
		dbNotificationSetting, dbNotificationFrequency, dbUser, dbFamily, dbZone, observer)
	getNotificationSettingsHandler := notificationSettingPres.NewGetNotificationSettingsHandler(
		getNotificationSettingsCommand, observer)
	getBatterySettingsCommand := notificationSettingApp.NewGetBatterySettingsHandler(dbNotificationSetting,
		dbNotificationFrequency, dbUser, dbFamily, dbZone, observer)
	getBatterySettingsHandler := notificationSettingPres.NewGetBatterySettingsHandler(getBatterySettingsCommand, observer)
	getNotificationsSendersByReceiverCommand := notificationSettingApp.NewGetNotificationsSendersByReceiverHandler(
		dbNotificationSetting, dbUser)
	getNotificationsSendersByReceiverHandler := notificationSettingPres.NewGetNotificationsSendersByReceiverHandler(
		getNotificationsSendersByReceiverCommand)
	changeNotificationSettingCommand := notificationSettingApp.NewChangeNotificationSettingHandler(
		dbNotificationSetting, dbNotificationFrequency, dbMembership, dbUser, dbFamily, dbRole, dbZone, observer)
	changeNotificationSettingHandler := notificationSettingPres.NewChangeNotificationSettingHandler(
		changeNotificationSettingCommand, observer)
	deleteNotificationSettingCommand := notificationSettingApp.NewDeleteNotificationSettingHandler(
		dbNotificationSetting, observer)
	deleteNotificationSettingHandler := notificationSettingPres.NewDeleteNotificationSettingHandler(
		deleteNotificationSettingCommand, observer)

	addTgLinkCommand := integration.NewAddTgLinkHandler(dbUser, dbIntegration, observer)
	addTgLinkHandler := integration2.NewAddTgLinkHandler(addTgLinkCommand, observer)
	addTgConnectionCommand := integration.NewAddTgConnectHandler(dbUser, dbIntegration, observer)
	addTgConnectionHandler := integration2.NewAddTgConnectHandler(addTgConnectionCommand, observer)
	//getNotificationLogsCommand := notification_log.NewGetNotificationLogsHandler(dbNotificationLog)
	//getNotificationLogsHandler := notification_log2.NewGetNotificationLogsHandler(getNotificationLogsCommand)
	//getUserNotificationLogsCommand := notification_log.NewGetUserNotificationLogsHandler(
	//	dbNotificationLog, dbNotificationSetting)
	//getUserNotificationLogsHandler := notification_log2.NewGetUserNotificationLogsHandler(
	//	getUserNotificationLogsCommand)

	securityMiddleware := http.NewSecurityMiddleware(observer, hydraService)

	app.Post("/users", addUserHandler.Handle)                                                                      //
	app.Post("/tg_link", securityMiddleware.Handle, addTgLinkHandler.Handle)                                       //
	app.Post("/tg_connect", addTgConnectionHandler.Handle)                                                         //
	app.Get("/me", securityMiddleware.Handle, getMeHandler.Handle)                                                 //
	app.Put("/users/tracking", securityMiddleware.Handle, changeUserTrackingHandler.Handle)                        //
	app.Get("/families", securityMiddleware.Handle, getUserMembershipsHandler.Handle)                              // проверить
	app.Post("/families", securityMiddleware.Handle, addFamilyHandler.Handle)                                      //
	app.Patch("/families/:id", securityMiddleware.Handle, changeFamilyHandler.Handle)                              //
	app.Delete("/memberships/:id", securityMiddleware.Handle, deleteUserMembershipHandler.Handle)                  // доделать
	app.Patch("/memberships/:id", securityMiddleware.Handle, changeFamilyMembershipHandler.Handle)                 //
	app.Post("/invitations", securityMiddleware.Handle, addInvitationToFamilyHandler.Handle)                       //
	app.Get("/invitations/:id", getInvitationsToFamilyByReceiverIdHandler.Handle)                                  //
	app.Post("/invitations/:id/join", securityMiddleware.Handle, addUserMembershipHandler.Handle)                  // проверить
	app.Get("/families/:id/zones", securityMiddleware.Handle, getFamilyZonesHandler.Handle)                        //
	app.Post("/families/:id/zones", securityMiddleware.Handle, addFamilyZoneHandler.Handle)                        //
	app.Patch("/zones/:id", securityMiddleware.Handle, changeZoneHandler.Handle)                                   //
	app.Delete("/zones/:id", securityMiddleware.Handle, deleteZoneHandler.Handle)                                  //
	app.Get("/frequencies", securityMiddleware.Handle, getAllNotificationFrequenciesHandler.Handle)                //
	app.Post("/notifications_settings", securityMiddleware.Handle, addNotificationSettingHandler.Handle)           //
	app.Get("/zones/:id/notifications_settings", securityMiddleware.Handle, getNotificationSettingsHandler.Handle) // добавить больше детей и взрослых
	app.Get("/battery/notifications_settings", securityMiddleware.Handle, getBatterySettingsHandler.Handle)
	app.Post("/users/locations", securityMiddleware.Handle, addUserLocationsHandler.Handle)                    // проверить
	app.Get("/families/:id/senders", securityMiddleware.Handle, getAvailableNotificationSendersHandler.Handle) // проверить
	app.Get("/users", securityMiddleware.Handle, getUsersHandler.Handle)
	app.Get("/users/locations", securityMiddleware.Handle, getUserLocationsHandler.Handle)

	app.Patch("/notifications_settings", changeNotificationSettingHandler.Handle)
	app.Delete("/notifications_settings/:id", deleteNotificationSettingHandler.Handle)
	//app.Get("/users/:id", getUserHandler.Handle)
	//app.Delete("/families/:family_id/zones/:zone_id", deleteFamilyZoneHandler.Handle)
	//app.Get("/messenger_types", getAllMessengerTypesHandler.Handle)
	//app.Get("/notification_types", getAllNotificationTypesHandler.Handle)
	//app.Patch("/users", changeUserHandler.Handle)
	//app.Post("/zones", addZoneHandler.Handle)
	//app.Get("/users/:id/zones", getZonesByUserIDHandler.Handle)

	app.Patch("/users_restore/:id", restoreUserHandler.Handle)
	app.Get("/families/:id", getFamilyHandler.Handle)
	app.Get("/users_families/:user_id", getUserFamiliesHandler.Handle)

	app.Delete("/families/:id", deleteFamilyHandler.Handle)

	//app.Delete("/families_zones/:id", deleteFamilyZonesHandler.Handle)

	//app.Delete("/users/:id", deleteUserHandler.Handle)

	app.Get("/zones/:id", getZoneHandler.Handle)

	app.Get("/zones", getAllZonesHandler.Handle)

	//app.Post("/roles", addRoleHandler.Handle)
	app.Get("/roles/:id", securityMiddleware.Handle, getRoleHandler.Handle)
	app.Get("/roles", getAllRolesHandler.Handle)

	app.Post("/memberships_by_invitations/:invitation_id", addMembershipByInvitationHandler.Handle)

	app.Get("/memberships/:family_id", getMembershipsByFamilyHandler.Handle)
	app.Get("/users_invitations/:user_id", getInvitationsToFamilyByReceiverIdHandler.Handle)

	app.Get("/notifications_senders/:receiver_id", getNotificationsSendersByReceiverHandler.Handle)

	//app.Get("/notifications_logs", getNotificationLogsHandler.Handle)
	//app.Get("/users_notifications_logs/:user_id", getUserNotificationLogsHandler.Handle)

	log.Info("Starting app")
	//routes := app.GetRoutes()
	//for _, route := range routes {
	//	observer.Logger.Info().Msgf("Route added: [%s]: %s", route.Method, route.Path)
	//}
	err = app.Listen(fmt.Sprintf(":%s", cfg.API.Port))
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
