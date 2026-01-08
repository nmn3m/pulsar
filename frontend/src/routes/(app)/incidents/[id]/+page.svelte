<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api/client';
  import { wsStore } from '$lib/stores/websocket';
  import Button from '$lib/components/ui/Button.svelte';
  import type {
    IncidentWithDetails,
    UpdateIncidentRequest,
    IncidentSeverity,
    IncidentStatus,
    ResponderWithUser,
    TimelineEventWithUser,
    IncidentAlertWithDetails,
  } from '$lib/types/incident';
  import type { User } from '$lib/types/user';

  let incident: IncidentWithDetails | null = null;
  let isLoading = true;
  let error: string | null = null;

  // Edit incident state
  let showEditForm = false;
  let editTitle = '';
  let editDescription = '';
  let editSeverity: IncidentSeverity = 'medium';
  let editStatus: IncidentStatus = 'investigating';
  let editPriority = 'P3';

  // Add note state
  let showNoteForm = false;
  let noteText = '';

  // Add responder state
  let showResponderForm = false;
  let selectedUserId = '';
  let selectedResponderRole: 'incident_commander' | 'responder' = 'responder';
  let availableUsers: User[] = [];

  // Link alert state
  let showLinkAlertForm = false;
  let selectedAlertId = '';
  let availableAlerts: any[] = [];

  let unsubscribeWS: (() => void)[] = [];

  $: incidentId = $page.params.id;

  onMount(async () => {
    await loadIncident();
    await loadAvailableUsers();

    // Listen for WebSocket incident events
    unsubscribeWS.push(
      wsStore.on('incident.created', () => {
        loadIncident(); // Refresh incident data
      })
    );
    unsubscribeWS.push(
      wsStore.on('incident.updated', () => {
        loadIncident();
      })
    );
    unsubscribeWS.push(
      wsStore.on('incident.deleted', () => {
        // Redirect to incidents list if this incident was deleted
        goto('/incidents');
      })
    );
    unsubscribeWS.push(
      wsStore.on('incident.timeline_added', () => {
        loadIncident();
      })
    );
    unsubscribeWS.push(
      wsStore.on('incident.responder_added', () => {
        loadIncident();
      })
    );
    unsubscribeWS.push(
      wsStore.on('incident.responder_removed', () => {
        loadIncident();
      })
    );
    unsubscribeWS.push(
      wsStore.on('incident.alert_linked', () => {
        loadIncident();
      })
    );
    unsubscribeWS.push(
      wsStore.on('incident.alert_unlinked', () => {
        loadIncident();
      })
    );
  });

  onDestroy(() => {
    // Unsubscribe from WebSocket events
    unsubscribeWS.forEach((unsub) => unsub());
  });

  async function loadIncident() {
    isLoading = true;
    error = null;

    try {
      incident = await api.getIncident(incidentId);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to load incident';
    } finally {
      isLoading = false;
    }
  }

  async function loadAvailableUsers() {
    try {
      availableUsers = await api.listUsers();
    } catch (err) {
      console.error('Failed to load users:', err);
    }
  }

  function startEdit() {
    if (!incident) return;
    editTitle = incident.title;
    editDescription = incident.description || '';
    editSeverity = incident.severity;
    editStatus = incident.status;
    editPriority = incident.priority;
    showEditForm = true;
  }

  async function handleUpdate() {
    if (!incident) return;

    const data: UpdateIncidentRequest = {
      title: editTitle.trim(),
      description: editDescription.trim() || undefined,
      severity: editSeverity,
      status: editStatus,
      priority: editPriority as any,
    };

    try {
      await api.updateIncident(incidentId, data);
      await loadIncident();
      showEditForm = false;
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to update incident';
    }
  }

  async function handleDelete() {
    if (!confirm('Are you sure you want to delete this incident?')) return;

    try {
      await api.deleteIncident(incidentId);
      goto('/incidents');
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to delete incident';
    }
  }

  async function handleAddNote() {
    if (!noteText.trim()) return;

    try {
      await api.addIncidentNote(incidentId, { note: noteText.trim() });
      noteText = '';
      showNoteForm = false;
      await loadIncident();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to add note';
    }
  }

  async function handleAddResponder() {
    if (!selectedUserId) return;

    try {
      await api.addIncidentResponder(incidentId, {
        user_id: selectedUserId,
        role: selectedResponderRole,
      });
      selectedUserId = '';
      selectedResponderRole = 'responder';
      showResponderForm = false;
      await loadIncident();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to add responder';
    }
  }

  async function handleRemoveResponder(userId: string) {
    if (!confirm('Remove this responder?')) return;

    try {
      await api.removeIncidentResponder(incidentId, userId);
      await loadIncident();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to remove responder';
    }
  }

  async function handleLinkAlert() {
    if (!selectedAlertId) return;

    try {
      await api.linkAlertToIncident(incidentId, { alert_id: selectedAlertId });
      selectedAlertId = '';
      showLinkAlertForm = false;
      await loadIncident();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to link alert';
    }
  }

  async function handleUnlinkAlert(alertId: string) {
    if (!confirm('Unlink this alert from the incident?')) return;

    try {
      await api.unlinkAlertFromIncident(incidentId, alertId);
      await loadIncident();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to unlink alert';
    }
  }

  function getSeverityColor(severity: IncidentSeverity): string {
    switch (severity) {
      case 'critical':
        return 'bg-red-100 text-red-800';
      case 'high':
        return 'bg-orange-100 text-orange-800';
      case 'medium':
        return 'bg-yellow-100 text-yellow-800';
      case 'low':
        return 'bg-blue-100 text-blue-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  }

  function getStatusColor(status: IncidentStatus): string {
    switch (status) {
      case 'investigating':
        return 'bg-yellow-100 text-yellow-800';
      case 'identified':
        return 'bg-blue-100 text-blue-800';
      case 'monitoring':
        return 'bg-purple-100 text-purple-800';
      case 'resolved':
        return 'bg-green-100 text-green-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  }

  function formatTimestamp(dateStr: string): string {
    const date = new Date(dateStr);
    return date.toLocaleString();
  }

  function getEventIcon(eventType: string): string {
    switch (eventType) {
      case 'created':
        return '🆕';
      case 'status_changed':
        return '🔄';
      case 'severity_changed':
        return '⚠️';
      case 'responder_added':
        return '👤';
      case 'responder_removed':
        return '👋';
      case 'note_added':
        return '📝';
      case 'alert_linked':
        return '🔗';
      case 'alert_unlinked':
        return '🔓';
      case 'resolved':
        return '✅';
      default:
        return '📌';
    }
  }

  $: responderUsers = incident?.responders?.filter((r) => r.role === 'responder') || [];
  $: incidentCommanders =
    incident?.responders?.filter((r) => r.role === 'incident_commander') || [];
</script>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
  <!-- Back Button -->
  <div class="mb-4">
    <Button variant="secondary" on:click={() => goto('/incidents')}>← Back to Incidents</Button>
  </div>

  {#if isLoading}
    <div class="text-center py-12">
      <p class="text-gray-500">Loading incident...</p>
    </div>
  {:else if error && !incident}
    <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
      {error}
    </div>
  {:else if incident}
    <!-- Incident Header -->
    <div class="bg-white shadow rounded-lg p-6 mb-6">
      <div class="flex items-start justify-between mb-4">
        <div class="flex-1">
          <div class="flex items-center gap-2 mb-3">
            <span
              class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium {getSeverityColor(
                incident.severity
              )}"
            >
              {incident.severity.toUpperCase()}
            </span>
            <span
              class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium {getStatusColor(
                incident.status
              )}"
            >
              {incident.status.charAt(0).toUpperCase() + incident.status.slice(1)}
            </span>
            <span
              class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-gray-100 text-gray-800"
            >
              {incident.priority}
            </span>
          </div>
          <h1 class="text-2xl font-bold text-gray-900 mb-2">{incident.title}</h1>
          {#if incident.description}
            <p class="text-gray-600">{incident.description}</p>
          {/if}
          <div class="mt-4 text-sm text-gray-500">
            <p>Started: {formatTimestamp(incident.started_at)}</p>
            {#if incident.resolved_at}
              <p>Resolved: {formatTimestamp(incident.resolved_at)}</p>
            {/if}
          </div>
        </div>
        <div class="flex gap-2">
          <Button on:click={startEdit}>Edit</Button>
          <Button variant="danger" on:click={handleDelete}>Delete</Button>
        </div>
      </div>
    </div>

    <!-- Edit Form -->
    {#if showEditForm}
      <div class="bg-white shadow rounded-lg p-6 mb-6">
        <h2 class="text-xl font-semibold mb-4">Edit Incident</h2>
        <div class="space-y-4">
          <div>
            <label for="edit-title" class="block text-sm font-medium text-gray-700">Title</label>
            <input
              id="edit-title"
              type="text"
              bind:value={editTitle}
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm px-3 py-2 border"
            />
          </div>
          <div>
            <label for="edit-description" class="block text-sm font-medium text-gray-700"
              >Description</label
            >
            <textarea
              id="edit-description"
              bind:value={editDescription}
              rows="3"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm px-3 py-2 border"
            />
          </div>
          <div class="grid grid-cols-3 gap-4">
            <div>
              <label for="edit-severity" class="block text-sm font-medium text-gray-700"
                >Severity</label
              >
              <select
                id="edit-severity"
                bind:value={editSeverity}
                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm px-3 py-2 border"
              >
                <option value="critical">Critical</option>
                <option value="high">High</option>
                <option value="medium">Medium</option>
                <option value="low">Low</option>
              </select>
            </div>
            <div>
              <label for="edit-status" class="block text-sm font-medium text-gray-700">Status</label
              >
              <select
                id="edit-status"
                bind:value={editStatus}
                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm px-3 py-2 border"
              >
                <option value="investigating">Investigating</option>
                <option value="identified">Identified</option>
                <option value="monitoring">Monitoring</option>
                <option value="resolved">Resolved</option>
              </select>
            </div>
            <div>
              <label for="edit-priority" class="block text-sm font-medium text-gray-700"
                >Priority</label
              >
              <select
                id="edit-priority"
                bind:value={editPriority}
                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm px-3 py-2 border"
              >
                <option value="P1">P1 - Critical</option>
                <option value="P2">P2 - High</option>
                <option value="P3">P3 - Medium</option>
                <option value="P4">P4 - Low</option>
                <option value="P5">P5 - Info</option>
              </select>
            </div>
          </div>
          <div class="flex gap-2">
            <Button on:click={handleUpdate}>Save Changes</Button>
            <Button variant="secondary" on:click={() => (showEditForm = false)}>Cancel</Button>
          </div>
        </div>
      </div>
    {/if}

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Main Content -->
      <div class="lg:col-span-2 space-y-6">
        <!-- Timeline -->
        <div class="bg-white shadow rounded-lg p-6">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-xl font-semibold">Timeline</h2>
            <Button on:click={() => (showNoteForm = !showNoteForm)}>Add Note</Button>
          </div>

          {#if showNoteForm}
            <div class="mb-4 p-4 bg-gray-50 rounded-lg">
              <textarea
                bind:value={noteText}
                rows="3"
                placeholder="Add a note to the timeline..."
                class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm px-3 py-2 border"
              />
              <div class="mt-2 flex gap-2">
                <Button on:click={handleAddNote} disabled={!noteText.trim()}>Add Note</Button>
                <Button variant="secondary" on:click={() => (showNoteForm = false)}>Cancel</Button>
              </div>
            </div>
          {/if}

          <div class="flow-root">
            <ul class="-mb-8">
              {#each incident.timeline || [] as event, idx (event.id)}
                <li>
                  <div class="relative pb-8">
                    {#if idx !== (incident.timeline?.length || 0) - 1}
                      <span
                        class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200"
                        aria-hidden="true"
                      />
                    {/if}
                    <div class="relative flex space-x-3">
                      <div>
                        <span
                          class="h-8 w-8 rounded-full flex items-center justify-center ring-8 ring-white bg-gray-100"
                        >
                          {getEventIcon(event.event_type)}
                        </span>
                      </div>
                      <div class="flex min-w-0 flex-1 justify-between space-x-4 pt-1.5">
                        <div>
                          <p class="text-sm text-gray-900">{event.description}</p>
                          {#if event.user}
                            <p class="text-xs text-gray-500">by {event.user.full_name}</p>
                          {/if}
                        </div>
                        <div class="whitespace-nowrap text-right text-sm text-gray-500">
                          <time>{formatTimestamp(event.created_at)}</time>
                        </div>
                      </div>
                    </div>
                  </div>
                </li>
              {/each}
            </ul>
          </div>
        </div>

        <!-- Linked Alerts -->
        <div class="bg-white shadow rounded-lg p-6">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-xl font-semibold">Linked Alerts ({incident.alerts?.length || 0})</h2>
          </div>

          {#if incident.alerts && incident.alerts.length > 0}
            <ul class="divide-y divide-gray-200">
              {#each incident.alerts as { alert, linked_at } (alert.id)}
                <li class="py-3 flex items-center justify-between">
                  <div>
                    <p class="text-sm font-medium text-gray-900">{alert.message}</p>
                    <p class="text-xs text-gray-500">
                      Priority: {alert.priority} • Status: {alert.status}
                    </p>
                    <p class="text-xs text-gray-500">Linked: {formatTimestamp(linked_at)}</p>
                  </div>
                  <Button variant="danger" size="sm" on:click={() => handleUnlinkAlert(alert.id)}>
                    Unlink
                  </Button>
                </li>
              {/each}
            </ul>
          {:else}
            <p class="text-sm text-gray-500">No alerts linked to this incident yet.</p>
          {/if}
        </div>
      </div>

      <!-- Sidebar -->
      <div class="space-y-6">
        <!-- Incident Commanders -->
        <div class="bg-white shadow rounded-lg p-6">
          <h3 class="text-lg font-semibold mb-4">Incident Commanders</h3>
          {#if incidentCommanders.length > 0}
            <ul class="space-y-2">
              {#each incidentCommanders as responder (responder.id)}
                <li class="flex items-center justify-between">
                  <div>
                    <p class="text-sm font-medium text-gray-900">{responder.user.full_name}</p>
                    <p class="text-xs text-gray-500">{responder.user.email}</p>
                  </div>
                  <button
                    type="button"
                    on:click={() => handleRemoveResponder(responder.user_id)}
                    class="text-red-600 hover:text-red-800 text-xs"
                  >
                    Remove
                  </button>
                </li>
              {/each}
            </ul>
          {:else}
            <p class="text-sm text-gray-500">No incident commanders assigned</p>
          {/if}
        </div>

        <!-- Responders -->
        <div class="bg-white shadow rounded-lg p-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold">Responders</h3>
            <Button size="sm" on:click={() => (showResponderForm = !showResponderForm)}>Add</Button>
          </div>

          {#if showResponderForm}
            <div class="mb-4 p-3 bg-gray-50 rounded">
              <select
                bind:value={selectedUserId}
                class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm px-3 py-2 border mb-2"
              >
                <option value="">Select user...</option>
                {#each availableUsers as user (user.id)}
                  <option value={user.id}>{user.full_name} ({user.email})</option>
                {/each}
              </select>
              <select
                bind:value={selectedResponderRole}
                class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm px-3 py-2 border mb-2"
              >
                <option value="incident_commander">Incident Commander</option>
                <option value="responder">Responder</option>
              </select>
              <div class="flex gap-2">
                <Button size="sm" on:click={handleAddResponder} disabled={!selectedUserId}
                  >Add</Button
                >
                <Button size="sm" variant="secondary" on:click={() => (showResponderForm = false)}
                  >Cancel</Button
                >
              </div>
            </div>
          {/if}

          {#if responderUsers.length > 0}
            <ul class="space-y-2">
              {#each responderUsers as responder (responder.id)}
                <li class="flex items-center justify-between">
                  <div>
                    <p class="text-sm font-medium text-gray-900">{responder.user.full_name}</p>
                    <p class="text-xs text-gray-500">{responder.user.email}</p>
                  </div>
                  <button
                    type="button"
                    on:click={() => handleRemoveResponder(responder.user_id)}
                    class="text-red-600 hover:text-red-800 text-xs"
                  >
                    Remove
                  </button>
                </li>
              {/each}
            </ul>
          {:else}
            <p class="text-sm text-gray-500">No responders assigned</p>
          {/if}
        </div>
      </div>
    </div>
  {/if}
</div>
