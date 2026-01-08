import { writable } from 'svelte/store';
import { browser } from '$app/environment';

export type WSEventType =
  | 'alert.created'
  | 'alert.updated'
  | 'alert.deleted'
  | 'alert.acknowledged'
  | 'alert.closed'
  | 'alert.escalated'
  | 'incident.created'
  | 'incident.updated'
  | 'incident.deleted'
  | 'incident.timeline_added'
  | 'incident.responder_added'
  | 'incident.responder_removed'
  | 'incident.alert_linked'
  | 'incident.alert_unlinked'
  | 'connection.connected'
  | 'connection.error'
  | 'connection.ping'
  | 'connection.pong';

export interface WSMessage {
  id: string;
  type: WSEventType;
  organization_id: string;
  payload: Record<string, any>;
  timestamp: string;
}

export type WSConnectionStatus = 'disconnected' | 'connecting' | 'connected' | 'error';

export type WSEventHandler = (message: WSMessage) => void;

interface WebSocketState {
  status: WSConnectionStatus;
  error: string | null;
  lastMessage: WSMessage | null;
}

function createWebSocketStore() {
  const { subscribe, set, update } = writable<WebSocketState>({
    status: 'disconnected',
    error: null,
    lastMessage: null,
  });

  let socket: WebSocket | null = null;
  let reconnectTimeout: ReturnType<typeof setTimeout> | null = null;
  let reconnectAttempts = 0;
  const maxReconnectAttempts = 5;
  const reconnectDelay = 3000; // 3 seconds

  const eventHandlers: Map<WSEventType, Set<WSEventHandler>> = new Map();

  function getWebSocketURL(): string {
    if (!browser) return '';

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.hostname;
    const port = import.meta.env.VITE_API_PORT || '8080';

    // Get token from localStorage
    const token = localStorage.getItem('token');

    return `${protocol}//${host}:${port}/api/v1/ws?token=${token}`;
  }

  function connect() {
    if (!browser) return;

    const token = localStorage.getItem('token');
    if (!token) {
      update((state) => ({ ...state, status: 'error', error: 'No authentication token' }));
      return;
    }

    if (socket?.readyState === WebSocket.OPEN) {
      return; // Already connected
    }

    update((state) => ({ ...state, status: 'connecting', error: null }));

    try {
      socket = new WebSocket(getWebSocketURL());

      socket.onopen = () => {
        console.log('WebSocket connected');
        reconnectAttempts = 0;
        update((state) => ({ ...state, status: 'connected', error: null }));
      };

      socket.onmessage = (event) => {
        try {
          const message: WSMessage = JSON.parse(event.data);

          // Update last message
          update((state) => ({ ...state, lastMessage: message }));

          // Call registered event handlers
          const handlers = eventHandlers.get(message.type);
          if (handlers) {
            handlers.forEach((handler) => handler(message));
          }

          // Call wildcard handlers (listening to all events)
          const wildcardHandlers = eventHandlers.get('*' as WSEventType);
          if (wildcardHandlers) {
            wildcardHandlers.forEach((handler) => handler(message));
          }
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };

      socket.onerror = (error) => {
        console.error('WebSocket error:', error);
        update((state) => ({
          ...state,
          status: 'error',
          error: 'WebSocket connection error',
        }));
      };

      socket.onclose = () => {
        console.log('WebSocket disconnected');
        update((state) => ({ ...state, status: 'disconnected' }));

        // Attempt to reconnect
        if (reconnectAttempts < maxReconnectAttempts) {
          reconnectAttempts++;
          console.log(`Reconnecting... (attempt ${reconnectAttempts}/${maxReconnectAttempts})`);
          reconnectTimeout = setTimeout(() => {
            connect();
          }, reconnectDelay);
        } else {
          update((state) => ({
            ...state,
            status: 'error',
            error: 'Max reconnection attempts reached',
          }));
        }
      };
    } catch (error) {
      console.error('Failed to create WebSocket:', error);
      update((state) => ({
        ...state,
        status: 'error',
        error: 'Failed to create WebSocket connection',
      }));
    }
  }

  function disconnect() {
    if (reconnectTimeout) {
      clearTimeout(reconnectTimeout);
      reconnectTimeout = null;
    }

    if (socket) {
      socket.close();
      socket = null;
    }

    reconnectAttempts = 0;
    update((state) => ({ ...state, status: 'disconnected', error: null }));
  }

  function on(eventType: WSEventType | '*', handler: WSEventHandler) {
    if (!eventHandlers.has(eventType as WSEventType)) {
      eventHandlers.set(eventType as WSEventType, new Set());
    }
    eventHandlers.get(eventType as WSEventType)!.add(handler);

    // Return unsubscribe function
    return () => {
      const handlers = eventHandlers.get(eventType as WSEventType);
      if (handlers) {
        handlers.delete(handler);
        if (handlers.size === 0) {
          eventHandlers.delete(eventType as WSEventType);
        }
      }
    };
  }

  function send(message: any) {
    if (socket?.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify(message));
    } else {
      console.warn('WebSocket is not connected');
    }
  }

  return {
    subscribe,
    connect,
    disconnect,
    on,
    send,
  };
}

export const wsStore = createWebSocketStore();
