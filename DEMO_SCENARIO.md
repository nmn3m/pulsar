# Pulsar Demo Scenario
## Complete Incident Management Workflow Presentation

---

## Overview

This demo showcases Pulsar as an end-to-end incident management platform, demonstrating how alerts flow from detection to resolution through teams, schedules, escalation policies, and notifications.

**Demo Duration:** ~15-20 minutes

---

## Pre-Demo Setup

### 1. Seed Demo Data (Recommended)
```bash
# Set required environment variables
export DATABASE_URL='postgres://user:pass@localhost:5432/pulsar?sslmode=disable'
export JWT_SECRET='your-secret-key-at-least-32-characters'
export JWT_REFRESH_SECRET='your-refresh-secret-at-least-32-chars'

# Run the seed script
cd backend
./scripts/seed-demo.sh
```

This creates:
- **Organization:** ACME Corporation
- **6 Users:** admin, alice, bob, carol, david, emma (password: `DemoPass123!`)
- **4 Teams:** Platform Engineering, Backend Services, Frontend Team, On-Call Rotation
- **2 Schedules:** Platform On-Call (weekly), Backend On-Call (daily)
- **2 Escalation Policies:** Platform and Backend escalation chains
- **8 Alerts:** Various priorities and statuses
- **3 Incidents:** With timeline events and linked alerts

### 2. Start the Services
```bash
# Start backend
cd backend && go run cmd/api/main.go

# Start frontend (separate terminal)
cd frontend && npm run dev
```

### 3. Verify Setup
- Backend running at `localhost:8080`
- Frontend running at `localhost:5173`
- Login with: `admin@acme-corp.com` / `DemoPass123!`

---

## Demo Script

### Act 1: Login & Overview (2 min)

**Narrative:** "Let's start by logging into our Pulsar instance for ACME Corporation."

#### Step 1.1: Login
1. Navigate to `http://localhost:5173/login`
2. Enter credentials:
   - Email: `admin@acme-corp.com`
   - Password: `DemoPass123!`
3. Click "Login"
4. Arrive at the Dashboard

**Talking Point:** "Pulsar supports multi-tenant organizations. Each organization has isolated data and can have multiple teams. We've pre-configured ACME Corporation with teams, schedules, and sample data."

#### Step 1.2: Quick Dashboard Overview
1. Point out the metrics cards showing current state
2. Highlight that we have some open alerts and incidents already

**Talking Point:** "The dashboard gives you an immediate view of your operational health."

---

### Act 2: Team Configuration (2 min)

**Narrative:** "Let me show you how we organize our incident response teams."

#### Step 2.1: View Teams
1. Navigate to **Teams** in sidebar
2. Show the 4 pre-configured teams:
   - Platform Engineering (3 members)
   - Backend Services (3 members)
   - Frontend Team (2 members)
   - On-Call Rotation (4 members)

**Talking Point:** "Teams represent groups of people who respond to specific types of incidents."

#### Step 2.2: Team Details
1. Click on **"Platform Engineering"** team
2. Show the team members with their roles:
   - Alice Chen (Lead)
   - Bob Martinez (Member)
   - Carol Williams (Member)

**Talking Point:** "Team leads can manage the team configuration, while members participate in rotations and incident response."

#### Step 2.3: (Optional) Add New Member
1. Click **"Add Member"**
2. Show the two options:
   - Select existing user from dropdown
   - Invite by email (sends invitation)
3. Cancel or complete based on time

**Talking Point:** "You can easily add team members by selecting existing users or inviting new people by email."

---

### Act 3: On-Call Schedules (2 min)

**Narrative:** "Let me show you how we manage on-call rotations."

#### Step 3.1: View Schedules
1. Navigate to **Schedules** in sidebar
2. Show the 2 pre-configured schedules:
   - Platform On-Call (Weekly rotation)
   - Backend On-Call (Daily rotation)

**Talking Point:** "Schedules define who is on-call at any given time. We support daily, weekly, and custom rotation types."

#### Step 3.2: Schedule Details
1. Click on **"Platform On-Call"** schedule
2. Point out:
   - Current on-call user (calculated automatically)
   - Weekly rotation with 3 participants (Alice, Bob, Carol)
   - Handoff time: 9 AM on Mondays

**Talking Point:** "The system automatically calculates who is currently on-call based on the rotation configuration."

#### Step 3.3: Show Override Capability
1. Point out the **"Add Override"** button
2. Explain: "If someone needs to swap shifts, you can create an override for a specific time period."

**Talking Point:** "Overrides allow flexible shift swapping without modifying the base rotation."

---

### Act 4: Escalation Policies (2 min)

**Narrative:** "Escalation policies ensure that if the first responder doesn't acknowledge, we escalate to others."

#### Step 4.1: View Policies
1. Navigate to **Escalation Policies** in sidebar
2. Show the 2 pre-configured policies:
   - Platform Escalation (3 rules, repeats 2x)
   - Backend Escalation (2 rules, no repeat)

**Talking Point:** "Escalation policies define who gets notified and when if an alert isn't acknowledged."

#### Step 4.2: Policy Details
1. Click on **"Platform Escalation"** policy
2. Walk through the 3-tier escalation:
   - **Rule 1 (5 min):** Notify Platform On-Call schedule
   - **Rule 2 (10 min):** Notify entire Platform Engineering team
   - **Rule 3 (15 min):** Notify Admin user directly

**Talking Point:** "First, we notify whoever is currently on-call. If they don't acknowledge within 5 minutes, we notify the entire team. If still no response after 10 more minutes, we escalate to the admin."

#### Step 4.3: Explain Repeat
1. Point out the "Repeat 2 times" setting
2. Explain: "If all rules execute without acknowledgment, the entire cycle repeats twice before giving up."

**Talking Point:** "The repeat feature ensures critical alerts don't get lost even if the first cycle fails."

---

### Act 5: Notification Channels (2 min)

**Narrative:** "Let's configure how we want to be notified."

#### Step 5.1: View Channels
1. Navigate to **Notifications > Channels** in sidebar

#### Step 5.2: Create Slack Channel (if applicable)
1. Click **"Create Channel"**
2. Select Type: `Slack`
3. Fill form:
   - Name: `Platform Slack Alerts`
   - Webhook URL: `https://hooks.slack.com/services/...`
   - Channel: `#platform-alerts`
4. Click "Create"

**Talking Point:** "Pulsar supports multiple notification channels including Email, Slack, Microsoft Teams, and custom Webhooks."

---

### Act 6: The Incident Scenario (5 min)

**Narrative:** "Now let's see how alerts and incidents work in practice with our pre-loaded data."

#### Step 6.1: View Active Alerts
1. Navigate to **Alerts** in sidebar
2. Point out the different alerts:
   - **P1 (Critical):** High CPU usage - Open
   - **P1 (Critical):** Database connection pool - Acknowledged
   - **P2 (High):** Memory leak detected - Open
   - Various other alerts in different states

**Talking Point:** "In production, these alerts would come from monitoring tools like Prometheus, Datadog, or CloudWatch via our API or webhooks."

#### Step 6.2: Alert Details
1. Click on the **"High CPU usage"** alert
2. Point out:
   - Status: `Open`
   - Priority: `P1` badge
   - Source: Prometheus
   - Tags: production, api, performance
   - Escalation policy assigned

**Talking Point:** "Each alert has full context about what triggered it and which escalation policy is handling notifications."

#### Step 6.3: Acknowledge Alert (Live Demo)
1. Click **"Acknowledge"** button
2. Notice status changes to `Acknowledged`

**Talking Point:** "Acknowledging stops the escalation timer. The responder is now actively working on the issue."

#### Step 6.4: View Existing Incident
1. Navigate to **Incidents**
2. Click on **"Production API Performance Degradation"** incident
3. Show:
   - Severity: Critical
   - Status: Investigating
   - Linked alert (High CPU alert)
   - Responders with roles (Alice as Incident Commander)
   - Timeline with investigation notes

**Talking Point:** "Incidents group related alerts together and track the full response lifecycle."

#### Step 6.5: Add Timeline Note (Live Demo)
1. In the timeline section, type a note:
   - `Demo: Continuing investigation. Checking memory metrics.`
2. Click "Add Note"
3. Note appears in timeline with timestamp

**Talking Point:** "The timeline provides a complete audit trail for post-mortems and compliance."

#### Step 6.6: View Resolved Incident
1. Go back to incidents list
2. Click on **"Database Connection Issues"** (Resolved)
3. Show:
   - Full timeline from creation to resolution
   - Resolution time tracked

**Talking Point:** "Resolved incidents maintain their full history for analysis and reporting."

#### Step 6.7: (Optional) Create New Incident
1. Click **"Create Incident"**
2. Show the form fields
3. Demonstrate linking multiple alerts

**Talking Point:** "When a major issue occurs, you can quickly create an incident and link all related alerts."

---

### Act 7: Dashboard & Metrics (2 min)

**Narrative:** "Finally, let's look at our operational metrics."

#### Step 7.1: View Dashboard
1. Navigate to **Dashboard**
2. Point out:
   - Open Alerts count
   - Active Incidents count
   - Mean Time to Acknowledge (MTTA)
   - Mean Time to Resolve (MTTR)

#### Step 7.2: Change Time Period
1. Click "Weekly" to see weekly metrics
2. Click "Monthly" to see monthly trends

**Talking Point:** "The dashboard gives leadership visibility into operational health and team performance."

#### Step 7.3: Show Charts
- Priority breakdown chart
- Source breakdown chart
- Team performance metrics

---

## Bonus Features to Mention

### API Keys
- Navigate to **Settings > API Keys**
- Show how to create API keys for CI/CD integration
- Mention scope-based permissions

### Routing Rules
- Navigate to **Settings > Routing Rules**
- Explain how alerts can be automatically routed to teams based on conditions

### Do Not Disturb
- Navigate to **Settings > DND**
- Show how users can set quiet hours

### Webhooks
- Navigate to **Webhooks**
- Show incoming webhook tokens for integrations
- Show outgoing webhook endpoints for notifications

---

## Key Talking Points Summary

1. **Multi-tenant Architecture** - Organizations with isolated data
2. **Flexible Team Structure** - Teams with roles and permissions
3. **Smart Scheduling** - Rotations with automatic on-call calculation
4. **Escalation Policies** - Ensure no alert goes unnoticed
5. **Multi-channel Notifications** - Slack, Teams, Email, Webhooks
6. **Incident Management** - Full lifecycle with timeline and responders
7. **Alert Correlation** - Link multiple alerts to incidents
8. **Metrics & Analytics** - MTTA, MTTR, team performance
9. **API-First Design** - Full REST API for integrations
10. **Audit Trail** - Complete history for compliance and post-mortems

---

## Seed Data Reference

### Users
| Name | Email | Role | Password |
|------|-------|------|----------|
| Demo Administrator | admin@acme-corp.com | Owner | DemoPass123! |
| Alice Chen | alice@acme-corp.com | Admin | DemoPass123! |
| Bob Martinez | bob@acme-corp.com | Member | DemoPass123! |
| Carol Williams | carol@acme-corp.com | Member | DemoPass123! |
| David Kim | david@acme-corp.com | Member | DemoPass123! |
| Emma Johnson | emma@acme-corp.com | Admin | DemoPass123! |

### Teams
| Team | Members | Description |
|------|---------|-------------|
| Platform Engineering | Alice (Lead), Bob, Carol | Infrastructure and platform reliability |
| Backend Services | Emma (Lead), David, Bob | API and microservices |
| Frontend Team | Carol (Lead), Alice | Web and mobile UI |
| On-Call Rotation | Admin (Lead), Alice, Emma, Bob | Cross-functional incident response |

### Schedules
| Schedule | Type | Participants | Handoff |
|----------|------|--------------|---------|
| Platform On-Call | Weekly | Alice, Bob, Carol | Monday 9 AM |
| Backend On-Call | Daily | Emma, David | 9 AM daily |

### Escalation Policies
| Policy | Rules | Repeat |
|--------|-------|--------|
| Platform Escalation | 1. On-Call (5m) → 2. Team (10m) → 3. Admin (15m) | 2x |
| Backend Escalation | 1. On-Call (5m) → 2. Emma (10m) | No |

### Sample Alerts
| Message | Priority | Status | Source |
|---------|----------|--------|--------|
| High CPU usage on prod-api-server-01 | P1 | Open | Prometheus |
| Database connection pool exhausted | P1 | Acknowledged | Datadog |
| Memory leak in payment-service | P2 | Open | NewRelic |
| SSL certificate expiring in 7 days | P3 | Open | CertManager |
| High error rate on checkout endpoint | P2 | Closed | Prometheus |
| Disk space warning on log-aggregator | P4 | Open | CloudWatch |
| K8s pod crashlooping: notification-worker | P2 | Open | Kubernetes |
| Unusual login activity detected | P3 | Acknowledged | SecurityHub |

### Sample Incidents
| Title | Severity | Status |
|-------|----------|--------|
| Production API Performance Degradation | Critical | Investigating |
| Payment Service Memory Issues | High | Identified |
| Database Connection Issues | Critical | Resolved |

---

## Demo Reset Commands

```bash
# Re-run the seed script to reset demo data
cd backend
./scripts/seed-demo.sh

# Or manually clear and reseed
psql $DATABASE_URL -c "TRUNCATE users, organizations, teams, schedules, alerts, incidents CASCADE;"
./scripts/seed-demo.sh
```

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Can't login | Check backend is running, verify credentials |
| No teams showing | Ensure you're in the correct organization |
| Notifications not sending | Check SMTP config for email, webhook URLs for Slack |
| Schedule not showing on-call | Ensure rotation has participants and valid start time |

---

## Architecture Diagram (for slides)

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

**End of Demo Script**
