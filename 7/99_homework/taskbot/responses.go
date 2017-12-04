package main

const (
	taskFormatResponse             = `%d. %s by @%s`
	taskCreatedResponse            = `Задача "%s" создана, id=%d`
	taskAssignedToYouResponse      = `Задача "%s" назначена на вас`
	taskAssignedToUserResponse     = `Задача "%s" назначена на @%s`
	taskUnassignAcceptedResponse   = `Принято`
	taskWithoutImplementerResponse = `Задача "%s" осталась без исполнителя`
	taskDoneResponse               = `Задача "%s" выполнена`
	taskDoneByResponse             = `Задача "%s" выполнена @%s`

	assigneeMe       = `assignee: я`
	assigneeUser     = `assignee: @%s`
	assignedFormat   = `assign_%d`
	unassignedFormat = `unassign_%d`
	resolveFormat    = `resolve_%d`
)
