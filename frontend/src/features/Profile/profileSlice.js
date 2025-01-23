import { createSlice } from '@reduxjs/toolkit';
import * as XLSX from 'xlsx';
import dayjs from 'dayjs';

const initialState = {
  loading: false,
  error: '',
  profiles: [],
  profileStats: {},
  maintenanceNotifications: [],
  profileDetails: {},
  avatar: '',
  recentActivities: [],
  favItems: [],
};

const profileSlice = createSlice({
  name: 'profile',
  initialState,
  reducers: {
    getProfileList: (state) => {
      state.error = '';
      state.profiles = [];
    },
    getProfileListSuccess: (state, action) => {
      state.profiles = action.payload;
      state.loading = false;
      state.error = '';
    },
    getProfileListFailure: (state) => {
      state.loading = false;
      state.error = '';
      state.profiles = [];
    },
    getProfileStats: (state) => {
      state.error = '';
      state.profileStats = {};
    },
    getProfileStatsSuccess: (state, action) => {
      state.profileStats = action.payload;
      state.loading = false;
      state.error = '';
    },
    getProfileStatsFailure: (state) => {
      state.loading = false;
      state.error = '';
      state.profileStats = {};
    },
    getMaintenanceNotifications: (state) => {
      state.error = '';
      state.maintenanceNotifications = [];
    },
    getMaintenanceNotificationsSuccess: (state, action) => {
      state.maintenanceNotifications = action.payload;
      state.loading = false;
      state.error = '';
    },
    getMaintenanceNotificationsFailure: (state) => {
      state.loading = false;
      state.error = '';
      state.maintenanceNotifications = [];
    },
    toggleMaintenanceNotificationReadOption: (state) => {
      state.error = '';
    },
    toggleMaintenanceNotificationReadOptionSuccess: (state, action) => {
      const { maintenance_plan_id, is_read, updated_at, updated_by } = action.payload;

      const index = state.maintenanceNotifications.findIndex(
        (item) => item.maintenance_plan_id === maintenance_plan_id
      );

      if (index !== -1) {
        state.maintenanceNotifications[index] = {
          ...state.maintenanceNotifications[index],
          is_read,
          updated_at,
          updated_by,
        };
      }
      state.error = '';
    },
    toggleMaintenanceNotificationReadOptionFailure: (state) => {
      state.error = '';
    },
    getRecentActivities: (state) => {
      state.loading = true;
      state.error = '';
    },
    getRecentActivitiesSuccess: (state, action) => {
      state.recentActivities = action.payload;
      state.loading = false;
      state.error = '';
    },
    getRecentActivitiesFailure: (state) => {
      state.loading = false;
      state.error = '';
      state.recentActivities = [];
    },
    getFavItems: (state) => {
      // does not meddle with profile loading
      state.loading = false;
      state.error = '';
    },
    getFavItemsSuccess: (state, action) => {
      state.favItems = action.payload;
      state.loading = false;
      state.error = '';
    },
    getFavItemsFailure: (state) => {
      state.loading = false;
      state.error = '';
      state.favItems = [];
    },
    saveFavItem: (state) => {
      // does not meddle with profile loading
      state.loading = false;
      state.error = '';
    },
    saveFavItemSuccess: (state, action) => {
      state.favItems = action.payload;
      state.loading = false;
      state.error = '';
    },
    saveFavItemFailure: (state) => {
      state.loading = false;
      state.error = '';
      state.favItems = [];
    },
    removeFavItem: (state) => {
      state.loading = true;
      state.error = '';
    },
    removeFavItemSuccess: (state, action) => {
      const filtered = [...state.favItems].filter((v) => {
        return v.id !== action.payload;
      });
      state.favItems = [...filtered];
      state.loading = false;
      state.error = '';
    },
    removeFavItemFailure: (state) => {
      state.loading = false;
      state.error = '';
      state.favItems = [];
    },
    downloadRecentActivities: (state) => {
      state.loading = true;
      state.error = '';
    },
    downloadRecentActivitiesSuccess: (state, action) => {
      const rawData = action.payload;
      const formattedData = rawData.map((v) => {
        // eslint-disable-next-line no-unused-vars
        const { id, activity_id, created_by, updated_by, sharable_groups, ...rest } = v;
        return rest;
      });

      const ws = XLSX.utils.json_to_sheet(formattedData);
      const wsColsWidth = Object.values(formattedData[0]).map((v) => ({ wch: v.length > 10 ? v.length : 10 }));
      ws['!cols'] = wsColsWidth;
      const wb = XLSX.utils.book_new();
      XLSX.utils.book_append_sheet(wb, ws, dayjs().format('YYYY-MM-DD'));
      XLSX.writeFile(wb, 'fleetwise.xlsx');
      state.loading = false;
      state.error = '';
    },
    downloadRecentActivitiesFailure: (state) => {
      state.loading = false;
      state.error = '';
    },
    getProfileDetails: (state) => {
      state.loading = true;
      state.error = '';
      state.profileDetails = {};
    },
    getProfileDetailsSuccess: (state, action) => {
      state.profileDetails = action.payload;
      state.loading = false;
      state.error = '';
    },
    getProfileDetailsFailure: (state) => {
      state.loading = false;
      state.error = '';
      state.profileDetails = {};
    },
    updateProfileDetails: (state) => {
      state.loading = true;
      state.error = '';
      state.profileDetails = {};
    },
    updateProfileDetailsSuccess: (state, action) => {
      state.profileDetails = action.payload;
      state.loading = false;
      state.error = '';
    },
    updateProfileDetailsFailure: (state) => {
      state.loading = false;
      state.error = '';
      state.profileDetails = {};
    },
    updateProfileImage: (state) => {
      state.loading = true;
      state.error = '';
    },
    updateProfileImageSuccess: (state) => {
      state.loading = false;
      state.error = '';
    },
    updateProfileImageFailure: (state) => {
      state.loading = false;
      state.error = '';
    },
    fetchAvatar: (state) => {
      state.error = '';
    },
    fetchAvatarSuccess: (state, action) => {
      state.loading = false;
      state.error = '';
      state.avatar = action.payload;
    },
    fetchAvatarFailure: (state) => {
      state.loading = false;
      state.error = '';
    },
  },
});

export const { actions: profileActions } = profileSlice;
export default profileSlice.reducer;
