package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
)

// AlertNotifier handles sending notifications for alert events
type AlertNotifier struct {
	notificationUsecase *NotificationUsecase
	userRepo            AlertNotifierUserRepository
	teamRepo            AlertNotifierTeamRepository
	scheduleUsecase     *ScheduleUsecase
	dndUsecase          *DNDUsecase
}

type AlertNotifierUserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	ListByTeam(ctx context.Context, teamID uuid.UUID) ([]*domain.User, error)
}

type AlertNotifierTeamRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Team, error)
}

func NewAlertNotifier(
	notificationUsecase *NotificationUsecase,
	userRepo AlertNotifierUserRepository,
	teamRepo AlertNotifierTeamRepository,
	scheduleUsecase *ScheduleUsecase,
	dndUsecase *DNDUsecase,
) *AlertNotifier {
	return &AlertNotifier{
		notificationUsecase: notificationUsecase,
		userRepo:            userRepo,
		teamRepo:            teamRepo,
		scheduleUsecase:     scheduleUsecase,
		dndUsecase:          dndUsecase,
	}
}

// NotifyAlertCreated sends notifications when a new alert is created
func (n *AlertNotifier) NotifyAlertCreated(ctx context.Context, alert *domain.Alert) error {
	// For now, this is a placeholder that can be expanded in future phases
	// Future implementation would:
	// 1. Determine who should be notified based on escalation policy
	// 2. Get active notification channels
	// 3. Send notifications through appropriate channels
	return nil
}

// NotifyAlertAcknowledged sends notifications when an alert is acknowledged
func (n *AlertNotifier) NotifyAlertAcknowledged(ctx context.Context, alert *domain.Alert, acknowledgedBy uuid.UUID) error {
	// Placeholder for future implementation
	return nil
}

// NotifyAlertClosed sends notifications when an alert is closed
func (n *AlertNotifier) NotifyAlertClosed(ctx context.Context, alert *domain.Alert, closedBy uuid.UUID, reason string) error {
	// Placeholder for future implementation
	return nil
}

// NotifyAlertEscalated sends notifications when an alert escalates
func (n *AlertNotifier) NotifyAlertEscalated(
	ctx context.Context,
	alert *domain.Alert,
	escalationRule *domain.EscalationRule,
	targets []domain.EscalationTarget,
) error {
	if n.notificationUsecase == nil {
		return nil // Notification usecase not configured
	}

	// Get all notification channels for the organization
	channels, err := n.notificationUsecase.ListChannels(ctx, alert.OrganizationID)
	if err != nil {
		return fmt.Errorf("failed to list notification channels: %w", err)
	}

	if len(channels) == 0 {
		// No channels configured, skip notifications
		return nil
	}

	// Build the notification message
	subject := fmt.Sprintf("[%s] Alert Escalated: %s", alert.Priority, alert.Message)
	message := fmt.Sprintf(
		"Alert ID: %s\nPriority: %s\nStatus: %s\nMessage: %s\n\nEscalation Level: %d\n\n%s",
		alert.ID,
		alert.Priority,
		alert.Status,
		alert.Message,
		alert.EscalationLevel,
		getDescriptionOrDefault(alert.Description),
	)

	// Send notifications to each target
	for _, target := range targets {
		recipients, err := n.resolveEscalationTarget(ctx, target)
		if err != nil {
			// Log error but continue with other targets
			continue
		}

		// Check if target has notification channel override
		targetChannelConfig, _ := target.ParseNotificationChannels()
		var targetChannelTypes []string
		if targetChannelConfig != nil && len(targetChannelConfig.Channels) > 0 {
			targetChannelTypes = targetChannelConfig.Channels
		}

		for _, recipient := range recipients {
			// Check if user is in DND mode
			if n.dndUsecase != nil {
				inDND, err := n.dndUsecase.IsInDNDMode(ctx, recipient.UserID, alert.Priority)
				if err == nil && inDND {
					// User is in DND mode, skip notification
					continue
				}
			}

			// Send through appropriate channels
			for _, channel := range channels {
				if !channel.IsEnabled {
					continue
				}

				// If target has specific channel override, only use those channels
				if len(targetChannelTypes) > 0 {
					if !containsChannelType(targetChannelTypes, string(channel.ChannelType)) {
						continue
					}
				}

				recipientAddr := recipient.ContactInfo

				// Construct notification request
				req := &SendNotificationRequest{
					ChannelID: channel.ID,
					UserID:    &recipient.UserID,
					AlertID:   &alert.ID,
					Recipient: recipientAddr,
					Subject:   &subject,
					Message:   message,
				}

				// Send notification (errors are logged in the notification usecase)
				_, _ = n.notificationUsecase.SendNotification(ctx, alert.OrganizationID, req)
			}
		}
	}

	return nil
}

// RecipientInfo contains user contact information for notifications
type RecipientInfo struct {
	UserID      uuid.UUID
	ContactInfo string // email, slack user id, etc.
}

// resolveEscalationTarget resolves an escalation target to actual recipients
func (n *AlertNotifier) resolveEscalationTarget(
	ctx context.Context,
	target domain.EscalationTarget,
) ([]RecipientInfo, error) {
	var recipients []RecipientInfo

	switch target.TargetType {
	case domain.EscalationTargetTypeUser:
		user, err := n.userRepo.GetByID(ctx, target.TargetID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}

		recipients = append(recipients, RecipientInfo{
			UserID:      user.ID,
			ContactInfo: user.Email,
		})

	case domain.EscalationTargetTypeTeam:
		teamMembers, err := n.userRepo.ListByTeam(ctx, target.TargetID)
		if err != nil {
			return nil, fmt.Errorf("failed to list team members: %w", err)
		}

		for _, member := range teamMembers {
			recipients = append(recipients, RecipientInfo{
				UserID:      member.ID,
				ContactInfo: member.Email,
			})
		}

	case domain.EscalationTargetTypeSchedule:
		// Get on-call user for this schedule at current time
		if n.scheduleUsecase == nil {
			return nil, fmt.Errorf("schedule usecase not configured")
		}

		onCallUser, err := n.scheduleUsecase.GetOnCallUser(ctx, target.TargetID, time.Now())
		if err != nil {
			return nil, fmt.Errorf("failed to get on-call user: %w", err)
		}

		if onCallUser != nil {
			user, err := n.userRepo.GetByID(ctx, onCallUser.UserID)
			if err == nil {
				recipients = append(recipients, RecipientInfo{
					UserID:      user.ID,
					ContactInfo: user.Email,
				})
			}
		}
	}

	return recipients, nil
}

func getDescriptionOrDefault(description *string) string {
	if description != nil {
		return *description
	}
	return "No additional description provided."
}

// containsChannelType checks if a channel type is in the list of allowed types
func containsChannelType(channelTypes []string, channelType string) bool {
	for _, ct := range channelTypes {
		if ct == channelType {
			return true
		}
	}
	return false
}
