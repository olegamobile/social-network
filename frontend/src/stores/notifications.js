import { defineStore } from 'pinia';
import { ref } from 'vue';

// Attempt to get API URL from environment variables, with a fallback
const apiUrl = import.meta.env.VITE_API_URL || '/api';

export const useNotificationStore = defineStore('notifications', () => {
    // State
    const notifications = ref([]); // Stores the list of notification objects
    const unreadCount = ref(0);    // Stores the count of unread notifications

    // Actions

    /**
     * Fetches all notifications for the user from the backend.
     * Populates the notifications list and updates the unread count.
     */
    async function fetchNotifications() {
        try {
            const response = await fetch(`${apiUrl}/api/notifications`, { // Corrected API URL
                credentials: 'include', // Important for sending session cookies
            });
            if (!response.ok) {
                throw new Error(`Failed to fetch notifications: ${response.statusText}`);
            }
            const data = await response.json(); // Expects an array of notification objects

            notifications.value = data;

            // Calculate unread count from the fetched notifications
            if (notifications.value) {
                unreadCount.value = data.filter(n => !n.is_read).length; // Simplified unread count
            }
            //console.log('Notifications fetched and store updated:', notifications.value.length, 'total,', unreadCount.value, 'unread');

        } catch (error) {
            console.error('Error in fetchNotifications:', error);
            // Potentially set an error state here for UI to pick up
        }
    }

    /**
     * Adds a new notification to the beginning of the list and increments unread count.
     * Typically called when a new notification arrives via WebSocket.
     * @param {object} newNotification - The new notification object.
     */
    async function addNotification(newNotification) { // Made async
        if (!notifications.value) notifications.value = []

        notifications.value.unshift(newNotification);
        if (!newNotification.is_read) {
            unreadCount.value++;
        }
        //console.log('New notification added locally, unread count:', unreadCount.value);

        try {
            //console.log('Refreshing all notifications to update pending statuses...');
            await fetchNotifications(); // Call fetchNotifications to refresh all
            //console.log('Notifications refreshed after new one arrived.');
        } catch (error) {
            console.error('Error refreshing notifications after new one arrived:', error);
            // Decide if any specific error handling is needed here for the fetchNotifications failure
        }
    }

    /**
     * Replaces the entire list of notifications.
     * @param {array} newNotificationsArray - The new array of notification objects.
     */
    function setNotifications(newNotificationsArray) {
        notifications.value = newNotificationsArray;
    }

    /**
     * Sets the unread notification count directly.
     * @param {number} count - The new unread count.
     */
    function setUnreadCount(count) {
        unreadCount.value = count;
    }

    /**
     * Marks a specific notification as read on the backend and updates its state in the store.
     * @param {number} notificationId - The ID of the notification to mark as read.
     */
    async function markAsRead(notificationId) {
        //console.log("marking notification", notificationId, "as read")

        try {
            const response = await fetch(`${apiUrl}/api/notifications/${notificationId}/read`, { // Corrected API URL
                method: 'POST',
                credentials: 'include',
            });

            if (!response.ok) {
                console.log("marking notification", notificationId, "as read failed miserably")
                throw new Error(`Failed to mark notification ${notificationId} as read: ${response.statusText}`);
            }

            // Update the local store
            const notification = notifications.value.find(n => n.id === notificationId);
            if (notification && !notification.is_read) {
                notification.is_read = true;
                if (unreadCount.value > 0) { // only decrement if it makes sense
                    unreadCount.value--;
                }
            }
            //console.log(`Notification ${notificationId} marked as read. Unread count:`, unreadCount.value);

        } catch (error) {
            console.log("Notification", notificationId, "as read failed at first hurdle:", error)
            console.error(`Error in markAsRead for notification ${notificationId}:`, error);
        }
    }

    // Return state and actions
    return {
        notifications,
        unreadCount,
        fetchNotifications,
        addNotification,
        setNotifications,
        setUnreadCount,
        markAsRead,
        updateNotification, // Added
        removeNotification, // Added
    };

    function updateNotification(updatedNotificationData) {
        const index = notifications.value.findIndex(n => n.id === updatedNotificationData.id);
        if (index !== -1) {
            notifications.value[index] = updatedNotificationData;
            unreadCount.value = notifications.value.filter(n => !n.is_read).length;
            console.log('Notification updated:', updatedNotificationData.id, 'Unread count:', unreadCount.value);
        } else {
            console.warn('Notification to update not found:', updatedNotificationData.id);
            // Optionally, if an update for a non-existent notification should add it:
            // addNotification(updatedNotificationData); 
            // However, current requirement is to log a warning.
        }
    }

    function removeNotification(notificationId) {
        const initialLength = notifications.value.length;
        notifications.value = notifications.value.filter(n => n.id !== notificationId);
        if (notifications.value.length < initialLength) {
            unreadCount.value = notifications.value.filter(n => !n.is_read).length;
            console.log('Notification removed:', notificationId, 'Unread count:', unreadCount.value);
        } else {
            console.warn('Notification to remove not found:', notificationId);
        }
    }
});
