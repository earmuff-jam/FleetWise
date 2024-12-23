import { takeLatest, put, call, takeEvery } from 'redux-saga/effects';
import { REACT_APP_LOCALHOST_URL } from '../../utils/Common';
import { profileActions } from './profileSlice';
import instance from '../../utils/Instances';

const DEFAULT_LIMIT = 10;
const BASEURL = `${REACT_APP_LOCALHOST_URL}/api/v1`;

export function* fetchProfileList() {
  try {
    const response = yield call(instance.get, `${BASEURL}/profile/list`);
    yield put(profileActions.getProfileListSuccess(response.data));
  } catch (e) {
    yield put(profileActions.getProfileListFailure(e));
  }
}

export function* fetchExistingUserDetails() {
  try {
    const USER_ID = localStorage.getItem('userID');
    const response = yield call(instance.get, `${BASEURL}/profile/${USER_ID}`);
    yield put(profileActions.getProfileDetailsSuccess(response.data));
  } catch (e) {
    yield put(profileActions.getProfileDetailsFailure(e));
  }
}

export function* fetchProfileStats() {
  try {
    const USER_ID = localStorage.getItem('userID');
    const response = yield call(instance.get, `${BASEURL}/profile/${USER_ID}/stats`);
    yield put(profileActions.getProfileStatsSuccess(response.data));
  } catch (e) {
    yield put(profileActions.getProfileStatsFailure(e));
  }
}

export function* fetchMaintenanceNotifications() {
  try {
    const USER_ID = localStorage.getItem('userID');
    const response = yield call(instance.get, `${BASEURL}/profile/${USER_ID}/notifications`);
    yield put(profileActions.getMaintenanceNotificationsSuccess(response.data));
  } catch (e) {
    yield put(profileActions.getMaintenanceNotificationsFailure(e));
  }
}

export function* fetchToggleMaintenanceNotificationReadOption(action) {
  try {
    const USER_ID = localStorage.getItem('userID');
    yield call(instance.put, `${BASEURL}/profile/${USER_ID}/notifications`, {
      ...action.payload,
      updated_by: USER_ID,
    });

    yield put(profileActions.toggleMaintenanceNotificationReadOptionSuccess(action.payload));
  } catch (e) {
    yield put(profileActions.toggleMaintenanceNotificationReadOptionFailure(e));
  }
}

export function* fetchRecentActivities() {
  try {
    const USER_ID = localStorage.getItem('userID');
    const params = new URLSearchParams();
    params.append('limit', DEFAULT_LIMIT);
    const response = yield call(instance.get, `${BASEURL}/profile/${USER_ID}/recent-activities?${params.toString()}`);
    yield put(profileActions.getRecentActivitiesSuccess(response.data));
  } catch (e) {
    yield put(profileActions.getRecentActivitiesFailure(e));
  }
}

export function* downloadRecentActivities(action) {
  try {
    const USER_ID = localStorage.getItem('userID');
    const params = new URLSearchParams();
    const { last30Days } = action.payload;
    params.append('until', last30Days);
    params.append('limit', 1000);
    const response = yield call(instance.get, `${BASEURL}/profile/${USER_ID}/recent-activities?${params.toString()}`);
    yield put(profileActions.downloadRecentActivitiesSuccess(response.data));
  } catch (e) {
    yield put(profileActions.downloadRecentActivitiesFailure(e));
  }
}

export function* updateExistingUserDetails(action) {
  try {
    const USER_ID = localStorage.getItem('userID');
    const { draftData } = action.payload;
    const response = yield call(instance.put, `${BASEURL}/profile/${USER_ID}`, {
      id: USER_ID,
      ...draftData,
      updated_by: USER_ID,
    });
    yield put(profileActions.getProfileDetailsSuccess(response.data));
  } catch (e) {
    yield put(profileActions.getProfileDetailsFailure(e));
  }
}

export function* fetchUpdateProfileImage(action) {
  try {
    const { id, selectedImage } = action.payload;
    const formData = new FormData();
    formData.append('imageSrc', selectedImage);
    const response = yield call(instance.post, `${BASEURL}/${id}/uploadImage`, formData);
    yield put(profileActions.updateProfileImageSuccess(response.data));
  } catch (e) {
    yield put(profileActions.updateProfileImageFailure(e));
  }
}

export function* fetchAvatar() {
  try {
    const USER_ID = localStorage.getItem('userID');
    // we need to modify the image to be of arrayBuffer type and build a blob object from it
    const response = yield call(instance.get, `${BASEURL}/${USER_ID}/fetchImage`, {
      responseType: 'arraybuffer',
    });
    const blob = new Blob([response.data], { type: response.headers['content-type'] });
    const avatarUrl = URL.createObjectURL(blob);
    yield put(profileActions.fetchAvatarSuccess(avatarUrl));
  } catch (e) {
    yield put(profileActions.fetchAvatarFailure(e));
  }
}

export function* fetchFavItems(action) {
  try {
    const USER_ID = localStorage.getItem('userID');
    const { limit } = action.payload;
    const params = new URLSearchParams();
    params.append('limit', limit);
    const response = yield call(instance.get, `${BASEURL}/profile/${USER_ID}/fav?${params.toString()}`);
    yield put(profileActions.getFavItemsSuccess(response.data));
  } catch (e) {
    yield put(profileActions.getFavItemsFailure(e));
  }
}

export function* saveFavItem(action) {
  try {
    const USER_ID = localStorage.getItem('userID');
    const draftPayload = { ...action.payload, created_by: USER_ID };
    const response = yield call(instance.post, `${BASEURL}/profile/${USER_ID}/fav`, { ...draftPayload });
    yield put(profileActions.saveFavItemSuccess(response.data));
  } catch (e) {
    yield put(profileActions.saveFavItemFailure(e));
  }
}

export function* removeFavItem(action) {
  try {
    const USER_ID = localStorage.getItem('userID');
    const params = new URLSearchParams();
    params.append('itemID', action.payload);
    const response = yield call(instance.delete, `${BASEURL}/profile/${USER_ID}/fav?${params.toString()}`);
    yield put(profileActions.removeFavItemSuccess(response.data));
  } catch (e) {
    yield put(profileActions.removeFavItemFailure(e));
  }
}

/********************************
 * WATCHER FUNCTIONS BELOW HERE
 ********************************/

export function* watchFetchProfileList() {
  yield takeLatest(`profile/getProfileList`, fetchProfileList);
}

export function* watchFetchExistingUserDetails() {
  yield takeLatest(`profile/getProfileDetails`, fetchExistingUserDetails);
}

export function* watchFetchProfileStats() {
  yield takeLatest(`profile/getProfileStats`, fetchProfileStats);
}

export function* watchFetchMaintenanceNotifications() {
  yield takeLatest(`profile/getMaintenanceNotifications`, fetchMaintenanceNotifications);
}

export function* watchFetchToggleMaintenanceNotificationReadOption() {
  yield takeLatest(`profile/toggleMaintenanceNotificationReadOption`, fetchToggleMaintenanceNotificationReadOption);
}

export function* watchFetchRecentActivities() {
  yield takeLatest(`profile/getRecentActivities`, fetchRecentActivities);
}

export function* watchDownloadRecentActivities() {
  yield takeLatest(`profile/downloadRecentActivities`, downloadRecentActivities);
}

export function* watchUpdateExistingUserDetails() {
  yield takeLatest(`profile/updateProfileDetails`, updateExistingUserDetails);
}

export function* watchFetchUpdateProfileImage() {
  yield takeEvery(`profile/updateProfileImage`, fetchUpdateProfileImage);
}

export function* watchFetchAvatar() {
  yield takeEvery(`profile/fetchAvatar`, fetchAvatar);
}

export function* watchFetchFavItems() {
  yield takeEvery(`profile/getFavItems`, fetchFavItems);
}

export function* watchSaveFavItem() {
  yield takeEvery(`profile/saveFavItem`, saveFavItem);
}

export function* watchRemoveFavItem() {
  yield takeEvery(`profile/removeFavItem`, removeFavItem);
}

export default [
  watchFetchAvatar,
  watchFetchProfileList,
  watchFetchProfileStats,
  watchFetchFavItems,
  watchSaveFavItem,
  watchRemoveFavItem,
  watchFetchRecentActivities,
  watchDownloadRecentActivities,
  watchFetchExistingUserDetails,
  watchUpdateExistingUserDetails,
  watchFetchUpdateProfileImage,
  watchFetchMaintenanceNotifications,
  watchFetchToggleMaintenanceNotificationReadOption,
];
