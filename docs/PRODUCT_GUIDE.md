# Pulsar - Incident Management Platform

## Product Overview

Pulsar is a comprehensive incident management platform designed for modern DevOps and SRE teams. It provides end-to-end alert management, on-call scheduling, escalation policies, and incident coordination to ensure your team never misses a critical issue.

---

## Table of Contents

1. [Authentication](#1-authentication)
2. [Dashboard](#2-dashboard)
3. [Alerts Management](#3-alerts-management)
4. [Incidents](#4-incidents)
5. [Teams](#5-teams)
6. [On-Call Schedules](#6-on-call-schedules)
7. [Escalation Policies](#7-escalation-policies)
8. [Notification Channels](#8-notification-channels)
9. [Webhooks](#9-webhooks)
10. [API Keys](#10-api-keys)
11. [Routing Rules](#11-routing-rules)
12. [Do Not Disturb](#12-do-not-disturb)

---

## 1. Authentication

![Login Page](screenshots/01-login.png)

### Features
- **Secure Login**: Email and password authentication with JWT tokens
- **Multi-tenant Support**: Each organization has isolated data
- **Session Management**: Automatic token refresh for seamless experience

### Demo Credentials
| User | Email | Password | Role |
|------|-------|----------|------|
| Admin | admin@acme-corp.com | DemoPass123! | Owner |
| Alice | alice@acme-corp.com | DemoPass123! | Admin |
| Bob | bob@acme-corp.com | DemoPass123! | Member |

---

## 2. Dashboard

![Dashboard](screenshots/02-dashboard.png)

### Features
- **Real-time Metrics**: View open alerts, active incidents, and key performance indicators
- **MTTA/MTTR Tracking**: Mean Time to Acknowledge and Mean Time to Resolve
- **Priority Distribution**: Visual breakdown of alerts by priority level
- **Source Analysis**: See which monitoring tools generate the most alerts
- **Team Performance**: Track response times across teams

### Key Metrics
- Open Alerts count
- Active Incidents count
- Alerts by Priority (P1-P5)
- Alerts by Source (Prometheus, Datadog, etc.)

---

## 3. Alerts Management

![Alerts Page](screenshots/03-alerts.png)

### Features
- **Priority Levels**: P1 (Critical) through P5 (Informational)
- **Status Tracking**: Open, Acknowledged, Snoozed, Closed
- **Source Integration**: Alerts from Prometheus, Datadog, CloudWatch, and more
- **Tagging System**: Organize alerts with custom tags
- **Quick Actions**: Acknowledge, Snooze, or Close alerts with one click

### Alert Statuses
| Status | Description |
|--------|-------------|
| Open | New alert requiring attention |
| Acknowledged | Responder is actively working on it |
| Snoozed | Temporarily silenced for a set duration |
| Closed | Issue resolved or dismissed |

### Sample Alerts
| Alert | Priority | Source | Status |
|-------|----------|--------|--------|
| High CPU usage on prod-api-server-01 | P1 | Prometheus | Open |
| Database connection pool exhausted | P1 | Datadog | Acknowledged |
| Memory leak in payment-service | P2 | NewRelic | Open |
| SSL certificate expiring in 7 days | P3 | CertManager | Open |
| Kubernetes pod crashlooping | P2 | Kubernetes | Open |

---

## 4. Incidents

![Incidents Page](screenshots/04-incidents.png)

### Features
- **Incident Lifecycle**: Track incidents from detection to resolution
- **Severity Levels**: Critical, High, Medium, Low
- **Responder Management**: Assign incident commanders and responders
- **Timeline Tracking**: Full audit trail of all actions taken
- **Alert Correlation**: Link multiple related alerts to a single incident

### Incident Statuses
| Status | Description |
|--------|-------------|
| Investigating | Initial triage and investigation |
| Identified | Root cause has been found |
| Monitoring | Fix deployed, monitoring for stability |
| Resolved | Issue fully resolved |

### Sample Incidents
| Incident | Severity | Status | Linked Alerts |
|----------|----------|--------|---------------|
| Production API Performance Degradation | Critical | Investigating | High CPU alert |
| Payment Service Memory Issues | High | Identified | Memory leak alert |
| Database Connection Issues | Critical | Resolved | DB connection alert |

---

## 5. Teams

![Teams Page](screenshots/05-teams.png)

### Features
- **Team Organization**: Group responders by function or service ownership
- **Role-based Access**: Team leads vs. members with different permissions
- **Member Management**: Add, remove, or update team members
- **Team Assignments**: Route alerts and incidents to specific teams

### Team Roles
| Role | Permissions |
|------|-------------|
| Lead | Manage team settings, add/remove members |
| Member | Respond to alerts and incidents |

### Sample Teams
| Team | Members | Description |
|------|---------|-------------|
| Platform Engineering | Alice (Lead), Bob, Carol | Infrastructure and platform reliability |
| Backend Services | Emma (Lead), David, Bob | API and microservices |
| Frontend Team | Carol (Lead), Alice | Web and mobile UI |
| On-Call Rotation | Admin (Lead), Alice, Emma, Bob | Cross-functional incident response |

---

## 6. On-Call Schedules

![Schedules Page](screenshots/06-schedules.png)

### Features
- **Flexible Rotations**: Daily, weekly, or custom rotation patterns
- **Multiple Schedules**: Different schedules for different teams/services
- **Override Support**: Temporary shift swaps without changing the rotation
- **Timezone Aware**: Proper handling of handoffs across timezones
- **Current On-Call**: Automatic calculation of who's currently on-call

### Rotation Types
| Type | Description |
|------|-------------|
| Daily | Rotate every day at handoff time |
| Weekly | Rotate every week on handoff day |
| Custom | Define your own rotation pattern |

### Sample Schedules
| Schedule | Type | Participants | Handoff |
|----------|------|--------------|---------|
| Platform On-Call | Weekly | Alice, Bob, Carol | Monday 9 AM |
| Backend On-Call | Daily | Emma, David | 9 AM daily |

---

## 7. Escalation Policies

![Escalation Policies Page](screenshots/07-escalations.png)

### Features
- **Multi-tier Escalation**: Define escalation chains with multiple levels
- **Flexible Targets**: Notify users, teams, or on-call schedules
- **Delay Configuration**: Set time between escalation levels
- **Repeat Cycles**: Optionally repeat the entire escalation chain
- **Policy Assignment**: Attach policies to alerts for automatic escalation

### Escalation Targets
| Target Type | Description |
|-------------|-------------|
| User | Notify a specific user |
| Team | Notify all team members |
| Schedule | Notify whoever is currently on-call |

### Sample Policies
| Policy | Rules | Repeat |
|--------|-------|--------|
| Platform Escalation | 1. On-Call (5m) → 2. Team (10m) → 3. Admin (15m) | 2x |
| Backend Escalation | 1. On-Call (5m) → 2. Emma (10m) | No |

---

## 8. Notification Channels

![Notification Channels](screenshots/08-notifications.png)

### Features
- **Multi-channel Support**: Email, Slack, Microsoft Teams, Webhooks
- **Channel Configuration**: Easy setup with validation
- **Enable/Disable**: Toggle channels without deleting configuration
- **Test Notifications**: Verify setup before going live

### Supported Channels
| Channel | Configuration |
|---------|---------------|
| Email | SMTP server, from address |
| Slack | Webhook URL, channel name |
| Microsoft Teams | Webhook URL |
| Webhook | Custom URL, headers, payload |

### Sample Channels
| Name | Type | Configuration |
|------|------|---------------|
| Team Email | Email | SMTP: mailpit:1025 |
| Slack #incidents | Slack | Channel: #incidents |
| Slack #alerts | Slack | Channel: #alerts |
| MS Teams Operations | Teams | Webhook configured |
| PagerDuty Integration | Webhook | events.pagerduty.com |

---

## 9. Webhooks

![Webhooks Page](screenshots/09-webhooks.png)

### Features
- **Outgoing Webhooks**: Send events to external systems
- **Incoming Webhooks**: Receive alerts from monitoring tools
- **Event Selection**: Choose which events trigger webhooks
- **Security**: HMAC signatures for payload verification
- **Retry Logic**: Automatic retries for failed deliveries

### Outgoing Webhook Events
| Event | Description |
|-------|-------------|
| alert.created | New alert created |
| alert.acknowledged | Alert acknowledged |
| alert.closed | Alert closed |
| incident.created | New incident created |
| incident.resolved | Incident resolved |

### Sample Outgoing Webhooks
| Name | URL | Events |
|------|-----|--------|
| DataDog Events | api.datadoghq.com | alert.created, incident.created |
| Jira Automation | atlassian.net | incident.created |
| Status Page | statuspage.io | incident.* |

### Incoming Webhook Integrations
| Integration | Default Priority | Use Case |
|-------------|------------------|----------|
| Prometheus | P2 | AlertManager integration |
| Grafana | P3 | Grafana alerting |
| Datadog | P2 | Datadog monitors |
| Generic | P3 | Custom integrations |

---

## 10. API Keys

![API Keys Page](screenshots/10-api-keys.png)

### Features
- **Scoped Access**: Fine-grained permissions per key
- **Expiration**: Set expiry dates for security
- **Key Management**: Create, revoke, and audit keys
- **Usage Tracking**: Monitor API key usage

### Available Scopes
| Scope | Permission |
|-------|------------|
| `*` | Full access to all endpoints |
| `alerts:read` | Read alerts |
| `alerts:write` | Create/update alerts |
| `incidents:read` | Read incidents |
| `incidents:write` | Create/update incidents |
| `teams:read` | Read team information |

### Sample API Keys
| Name | Owner | Scopes | Expiration |
|------|-------|--------|------------|
| Demo Integration Key | admin@acme-corp.com | * (full access) | 1 year |
| Read-Only Dashboard Key | alice@acme-corp.com | alerts:read, incidents:read | 1 year |

---

## 11. Routing Rules

![Routing Rules Page](screenshots/11-routing-rules.png)

### Features
- **Automatic Routing**: Route alerts based on conditions
- **Priority-based Execution**: Rules execute in priority order
- **Flexible Conditions**: Match on priority, source, tags, message
- **Rich Actions**: Assign teams, set priority, add tags, suppress

### Condition Operators
| Operator | Description |
|----------|-------------|
| equals | Exact match |
| contains | Substring match |
| matches | Regex pattern |
| exists | Field is present |

### Available Actions
| Action | Description |
|--------|-------------|
| Assign Team | Route to a specific team |
| Assign Policy | Apply an escalation policy |
| Set Priority | Override alert priority |
| Add Tags | Append tags to the alert |
| Suppress | Silently discard the alert |

### Sample Routing Rules
| Rule | Priority | Conditions | Actions |
|------|----------|------------|---------|
| Suppress Test Alerts | 0 | tags contains test/staging | Suppress |
| Critical to Platform | 1 | priority=P1 AND source=prometheus | Platform team |
| Database Alerts | 2 | tags contains database/postgres | Backend team |
| Security Boost | 3 | tags contains security | Set P2, add tag |

---

## 12. Do Not Disturb

![DND Settings](screenshots/12-dnd.png)

### Features
- **Personal Schedules**: Set quiet hours per user
- **Weekly Patterns**: Define recurring DND windows
- **P1 Override**: Critical alerts can bypass DND
- **Temporary Overrides**: One-time DND periods

### Configuration Options
| Option | Description |
|--------|-------------|
| Schedule | Weekly recurring quiet hours |
| Timezone | User's local timezone |
| P1 Override | Allow P1 alerts to bypass DND |
| Overrides | One-time DND periods |

### Sample DND Settings
| User | Schedule | P1 Override |
|------|----------|-------------|
| Alice Chen | Weeknights 10PM-8AM (Mon-Fri) | Enabled |
| Bob Martinez | Weekends (Sat-Sun all day) | Enabled |

---

## Architecture

```
                    +------------------+
                    |   Monitoring     |
                    | (Prometheus,     |
                    |  Datadog, etc)   |
                    +--------+---------+
                             |
                             | Webhooks/API
                             v
+------------------+  +------+-------+  +------------------+
|                  |  |              |  |                  |
|  Pulsar Frontend +--+   Pulsar    +--+ Notification     |
|   (SvelteKit)    |  |   Backend   |  | Channels         |
|                  |  |   (Go/Gin)  |  | (Slack/Email/    |
+------------------+  |              |  |  Teams/Webhook)  |
                      +------+-------+  +------------------+
                             |
                             v
                      +------+-------+
                      |  PostgreSQL  |
                      |   Database   |
                      +--------------+
```

---

## Quick Start

### 1. Start the Services
```bash
make up
```

### 2. Seed Demo Data
```bash
cd backend
DATABASE_URL="postgres://pulsar:pulsar_dev_password@localhost:5433/pulsar?sslmode=disable" \
JWT_SECRET="dev_jwt_secret_change_in_production_min_32_chars" \
JWT_REFRESH_SECRET="dev_refresh_secret_change_in_production_min_32_chars" \
go run ./cmd/seed/main.go
```

### 3. Access the Application
- Frontend: http://localhost:5173
- Backend API: http://localhost:8081
- API Docs: http://localhost:8081/swagger/index.html

### 4. Login
- Email: `admin@acme-corp.com`
- Password: `DemoPass123!`

---

## Key Benefits

1. **Unified Alert Management** - All your alerts in one place
2. **Smart Escalations** - Never miss a critical issue
3. **Flexible Scheduling** - On-call rotations that work for your team
4. **Multi-channel Notifications** - Reach responders wherever they are
5. **Full Audit Trail** - Complete history for compliance and post-mortems
6. **API-First Design** - Integrate with any tool in your stack
7. **Real-time Updates** - WebSocket-powered live updates

---

## Contact

For questions or support, please contact the Pulsar team.
