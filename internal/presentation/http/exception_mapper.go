package http

import (
	"awesomeProjectDDD/internal/adapter"
	"awesomeProjectDDD/internal/application"
	"awesomeProjectDDD/internal/domain"
	"awesomeProjectDDD/internal/presentation/http/response"
	"errors"
	"github.com/elliotchance/orderedmap/v3"
	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Status int `json:"status"`
	Data   response.BaseStatusErrorResponse
}

func NewHTTPError(status int, message string) Error {
	return Error{
		Status: status,
		Data:   response.NewErrorResponse(message),
	}
}

var ErrCodeMap = orderedmap.NewOrderedMap[error, Error]()

func init() {
	ErrCodeMap.Set(domain.ErrInvalidID, NewHTTPError(fiber.StatusUnprocessableEntity, "неверный id"))
	ErrCodeMap.Set(adapter.ErrUserNotFound, NewHTTPError(fiber.StatusNotFound, "пользователь не найден в базе"))
	ErrCodeMap.Set(domain.ErrInvalidPhone, NewHTTPError(fiber.StatusUnprocessableEntity, "неверный номер телефона"))
	ErrCodeMap.Set(application.ErrPhoneExists, NewHTTPError(fiber.StatusBadRequest, "номер телефона уже существует"))
	ErrCodeMap.Set(domain.ErrInvalidEmail, NewHTTPError(fiber.StatusUnprocessableEntity, "неверный адрес электронной почты"))
	ErrCodeMap.Set(application.ErrEmailExists, NewHTTPError(fiber.StatusBadRequest, "адрес электронной почты уже существует"))
	ErrCodeMap.Set(domain.ErrInvalidFirstName, NewHTTPError(fiber.StatusUnprocessableEntity, "неверное имя"))
	ErrCodeMap.Set(domain.ErrInvalidLastName, NewHTTPError(fiber.StatusUnprocessableEntity, "неверная фамилия"))
	ErrCodeMap.Set(domain.ErrInvalidLogin, NewHTTPError(fiber.StatusUnprocessableEntity, "неверный логин"))
	ErrCodeMap.Set(domain.ErrInvalidBirthDate, NewHTTPError(fiber.StatusUnprocessableEntity, "неверная дата рождения"))
	ErrCodeMap.Set(domain.ErrInvalidPassword, NewHTTPError(fiber.StatusUnprocessableEntity, "неверный пароль"))
	ErrCodeMap.Set(domain.ErrUserIsDeleted, NewHTTPError(fiber.StatusBadRequest, "пользователь удален"))
	ErrCodeMap.Set(adapter.ErrMembershipsNotFound, NewHTTPError(fiber.StatusNotFound, "участие в семье не найдено"))
	ErrCodeMap.Set(application.ErrMembershipUserNotFound, NewHTTPError(fiber.StatusBadRequest, "пользователь для участия в семье не найден"))
	ErrCodeMap.Set(application.ErrMembershipFamilyNotFound, NewHTTPError(fiber.StatusBadRequest, "семья для участия в семье не найдена"))
	ErrCodeMap.Set(adapter.ErrMembershipUsersNotFound, NewHTTPError(fiber.StatusNotFound, "пользователи-участники семьи не найдены в базе"))
	ErrCodeMap.Set(adapter.ErrMembershipFamiliesNotFound, NewHTTPError(fiber.StatusNotFound, "семьи для участий не найдены в базе"))
	ErrCodeMap.Set(domain.ErrInvalidFamilyName, NewHTTPError(fiber.StatusUnprocessableEntity, "неверное название семьи"))
	ErrCodeMap.Set(adapter.ErrFamilyNotFound, NewHTTPError(fiber.StatusNotFound, "семья не найдена"))
	ErrCodeMap.Set(application.ErrMembershipExists, NewHTTPError(fiber.StatusBadRequest, "участие в семье уже существует"))
	ErrCodeMap.Set(adapter.ErrRoleNotFound, NewHTTPError(fiber.StatusNotFound, "роль не найдена"))
	ErrCodeMap.Set(application.ErrUserIsNotAdult, NewHTTPError(fiber.StatusBadRequest, "пользователь не является родителем"))
	ErrCodeMap.Set(domain.ErrInvalidZoneName, NewHTTPError(fiber.StatusUnprocessableEntity, "неверое имя зоны"))
	ErrCodeMap.Set(domain.ErrInvalidZoneSafety, NewHTTPError(fiber.StatusUnprocessableEntity, "неверная безопасность зоны"))
	ErrCodeMap.Set(domain.ErrInvalidZoneBoundaries, NewHTTPError(fiber.StatusUnprocessableEntity, "неверные границы зоны"))
	ErrCodeMap.Set(adapter.ErrZoneNotFound, NewHTTPError(fiber.StatusNotFound, "зона не найдена"))
	ErrCodeMap.Set(adapter.ErrMembershipNotFound, NewHTTPError(fiber.StatusNotFound, "участие в семье не найдено"))
	ErrCodeMap.Set(adapter.ErrInvitationNotFound, NewHTTPError(fiber.StatusNotFound, "приглашение в семью не найдено"))
	ErrCodeMap.Set(adapter.ErrNotFamilyMembers, NewHTTPError(fiber.StatusBadRequest, "пользователи не члены одной семьи"))
	ErrCodeMap.Set(domain.ErrZoneIsDeleted, NewHTTPError(fiber.StatusBadRequest, "зона удалена"))
	ErrCodeMap.Set(domain.ErrInvalidSendTime, NewHTTPError(fiber.StatusUnprocessableEntity, "неверное время отправки"))
	ErrCodeMap.Set(domain.ErrInvalidNotificationLogContext, NewHTTPError(fiber.StatusUnprocessableEntity, "неверный контекст лога об отправленном уведомении"))
	ErrCodeMap.Set(domain.ErrInvalidFrequency, NewHTTPError(fiber.StatusUnprocessableEntity, "неверная частота отправки уведомлений"))
	ErrCodeMap.Set(domain.ErrInvalidEventType, NewHTTPError(fiber.StatusUnprocessableEntity, "неверный тип уведомления"))
	ErrCodeMap.Set(domain.ErrUserNotFound, NewHTTPError(fiber.StatusNotFound, "пользователь не найден (домен)"))
	ErrCodeMap.Set(domain.ErrInvalidZoneOrBattery, NewHTTPError(fiber.StatusUnprocessableEntity, "требуются зона или заряд батареи"))
	ErrCodeMap.Set(domain.ErrInvalidZone, NewHTTPError(fiber.StatusBadRequest, "неверная зона"))
	ErrCodeMap.Set(domain.ErrInvalidBattery, NewHTTPError(fiber.StatusBadRequest, "неверный заряд батареи"))
	ErrCodeMap.Set(domain.ErrNotifySettingIsDeleted, NewHTTPError(fiber.StatusBadRequest, "подписка на уведомление удалена"))
	ErrCodeMap.Set(domain.ErrFamilyIsDeleted, NewHTTPError(fiber.StatusBadRequest, "семья удалена"))
	ErrCodeMap.Set(domain.ErrRoleNotFound, NewHTTPError(fiber.StatusNotFound, "роль не найдена"))
	ErrCodeMap.Set(domain.ErrFamilyNotFound, NewHTTPError(fiber.StatusNotFound, "семья не найдена"))
	ErrCodeMap.Set(domain.ErrInvitationToFamilyIsDeleted, NewHTTPError(fiber.StatusBadRequest, "приглашение в семью удалено"))
	ErrCodeMap.Set(application.ErrReceiverIsSender, NewHTTPError(fiber.StatusBadRequest, "отправитель является получателем"))
	ErrCodeMap.Set(adapter.ErrFrequencyNotFound, NewHTTPError(fiber.StatusNotFound, "частота отправки не найдена"))
	ErrCodeMap.Set(adapter.ErrNotificationSettingExists, NewHTTPError(fiber.StatusBadRequest, "подписка на уведомление уже существует"))
	ErrCodeMap.Set(adapter.ErrNotificationSettingNotFound, NewHTTPError(fiber.StatusNotFound, "подписка на уведомление не найдена"))
	ErrCodeMap.Set(application.ErrDifferentArrayLength, NewHTTPError(fiber.StatusBadRequest, "разные длины массивов"))
	ErrCodeMap.Set(application.ErrNotifySettingsNotFound, NewHTTPError(fiber.StatusBadRequest, "подписки на уведомления не найдены"))
	ErrCodeMap.Set(application.ErrNotifyFrequenciesNotFound, NewHTTPError(fiber.StatusBadRequest, "частоты отправки уведомлений не найдены"))
	ErrCodeMap.Set(application.ErrNotifySendersNotFound, NewHTTPError(fiber.StatusBadRequest, "отправители уведомлений не найдены"))
	ErrCodeMap.Set(application.ErrNotifyZoneNotFound, NewHTTPError(fiber.StatusBadRequest, "зона для подписки на уведомления не найдена"))
	ErrCodeMap.Set(application.ErrInvalidNotifySettingDetail, NewHTTPError(fiber.StatusBadRequest, "неверные детали подписки на уведомления"))
	ErrCodeMap.Set(application.ErrZoneOverlapsZones, NewHTTPError(fiber.StatusBadRequest, "зона пересекается с другими зонами семьи"))
	ErrCodeMap.Set(application.ErrUserNotTracked, NewHTTPError(fiber.StatusBadRequest, "пользователь не отслеживается"))
	ErrCodeMap.Set(domain.ErrInvalidLatitude, NewHTTPError(fiber.StatusBadRequest, "неверная широта"))
	ErrCodeMap.Set(domain.ErrInvalidLongitude, NewHTTPError(fiber.StatusBadRequest, "неверная долгота"))
	ErrCodeMap.Set(domain.ErrUserAlreadyIntegrated, NewHTTPError(fiber.StatusBadRequest, "интеграция с системой уведомлений уже существует"))
	ErrCodeMap.Set(domain.ErrTgConnect, NewHTTPError(fiber.StatusBadRequest, "не удалось подключить систему уведомлений"))
	ErrCodeMap.Set(adapter.ErrIntegrationAlreadyExists, NewHTTPError(fiber.StatusBadRequest, "интеграция с системой уведомлений уже существует"))
	ErrCodeMap.Set(ErrAuthorizationNotFound, NewHTTPError(fiber.StatusUnauthorized, "не удалось авторизоваться"))
	ErrCodeMap.Set(ErrAccessTokenExpired, NewHTTPError(fiber.StatusUnauthorized, "токен доступа закончился"))

	ErrCodeMap.Set(domain.ErrDomain, NewHTTPError(fiber.StatusInternalServerError, "доменная ошибка"))
	ErrCodeMap.Set(application.ErrApplication, NewHTTPError(fiber.StatusInternalServerError, "ошибка бизнес-логики"))
	ErrCodeMap.Set(adapter.ErrStorage, NewHTTPError(fiber.StatusInternalServerError, "ошибка хранилища"))
	ErrCodeMap.Set(application.ErrGeneral, NewHTTPError(fiber.StatusInternalServerError, "ошибка приложения"))
}

func GetErrorCodeFromError(err error) (*Error, bool) {
	for k, v := range ErrCodeMap.AllFromFront() {
		if errors.Is(err, k) {
			return &v, true
		}
	}
	return nil, false
}

func RespondError(ctx *fiber.Ctx, err error) error {
	httpErr, ok := GetErrorCodeFromError(err)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			NewHTTPError(
				fiber.StatusInternalServerError,
				"internal server error",
			),
		)
	}

	return ctx.Status(httpErr.Status).JSON(httpErr.Data)
}
