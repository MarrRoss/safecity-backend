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
	"html/template"
	"time"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
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
		dbNotificationFrequency, dbNotificationSetting,
		dbUserLocations, dbNotificationLog, dbIntegration, err := db.CreateStorages(observer, pg)

	addUserCommand := user.NewAddUserHandler(dbUser, observer)
	addUserHandler := userPres.NewAddUserHandler(addUserCommand, observer)
	getUsersCommand := user.NewGetAllUsersHandler(dbUser, observer)
	getUsersHandler := userPres.NewGetAllUsersHandler(getUsersCommand, observer)
	getUserQuery := user.NewGetUserHandler(observer, dbUser)
	getMeHandler := userPres.NewGetMeHandler(getUserQuery, observer)
	changeUserTrackingCommand := user.NewChangeUserTrackingHandler(dbUser, observer)
	changeUserTrackingHandler := userPres.NewChangeUserTrackingHandler(changeUserTrackingCommand, observer)

	addFamilyCommand := family.NewAddFamilyHandler(dbFamily, dbUser, dbMembership, dbRole, observer)
	addFamilyHandler := familyPres.NewAddFamilyHandler(addFamilyCommand, observer)
	changeFamilyCommand := family.NewChangeFamilyHandler(dbFamily, observer)
	changeFamilyHandler := familyPres.NewChangeFamilyHandler(changeFamilyCommand, observer)

	changeZoneCommand := zone.NewChangeZoneHandler(dbZone, dbUser, observer)
	changeZoneHandler := zonePres.NewChangeZoneHandler(changeZoneCommand, observer)
	deleteZoneCommand := zone.NewDeleteZoneHandler(dbZone, dbUser, observer)
	deleteZoneHandler := zonePres.NewDeleteZoneHandler(deleteZoneCommand, observer)

	addFamilyZoneCommand := family.NewAddFamilyZoneHandler(
		dbFamily, dbZone, observer)
	addFamilyZoneHandler := familyPres.NewAddFamilyZoneHandler(addFamilyZoneCommand, observer)
	getFamilyZonesCommand := family.NewGetFamilyZonesHandler(dbFamily, observer)
	getFamilyZonesHandler := familyPres.NewGetFamilyZonesHandler(getFamilyZonesCommand, observer)

	getRoleCommand := role.NewGetRoleHandler(dbRole, hydraService)
	getRoleHandler := rolePres.NewGetRoleHandler(getRoleCommand)
	getAllRolesCommand := role.NewGetAllRolesHandler(dbRole)
	getAllRolesHandler := rolePres.NewGetAllRolesHandler(getAllRolesCommand)

	getUserMembershipsCommand := membershipApp.NewGetUserMembershipsHandler(dbMembership, dbUser, dbFamily, observer)
	getUserMembershipsHandler := membershipPres.NewGetUserMembershipsHandler(getUserMembershipsCommand, observer)
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
	app.Get("/roles/:id", securityMiddleware.Handle, getRoleHandler.Handle)
	app.Get("/roles", getAllRolesHandler.Handle)

	log.Info("Starting app")
	err = app.Listen(fmt.Sprintf(":%s", cfg.API.Port))
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
