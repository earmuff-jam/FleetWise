import { call, put, takeLatest } from 'redux-saga/effects';

import instance from '../../utils/Instances';
import { REACT_APP_LOCALHOST_URL } from '../../utils/Common';
import { maintenancePlanItemActions } from './maintenancePlanItemSlice';

const DEFAULT_LIMIT = 10;
const BASEURL = `${REACT_APP_LOCALHOST_URL}/api/v1`;

export function* getSelectedMaintenancePlan(action) {
  try {
    const mID = action.payload;
    const userID = localStorage.getItem('userID');
    const params = new URLSearchParams();
    params.append('id', userID);
    params.append('mID', mID);
    const response = yield call(instance.get, `${BASEURL}/plan?${params.toString()}`);
    yield put(maintenancePlanItemActions.getSelectedMaintenancePlanSuccess(response.data));
  } catch (e) {
    yield put(maintenancePlanItemActions.getSelectedMaintenancePlanFailure(e));
  }
}

export function* updateCollaborators(action) {
  try {
    const { id } = action.payload;
    const response = yield call(instance.put, `${BASEURL}/plan/${id}`, { ...action.payload });
    yield put(maintenancePlanItemActions.updateCollaboratorsSuccess(response.data));
  } catch (e) {
    yield put(maintenancePlanItemActions.updateCollaboratorsFailure(e));
  }
}

export function* getItemsInMaintenancePlan(action) {
  try {
    const mID = action.payload;
    const userID = localStorage.getItem('userID');
    const params = new URLSearchParams();
    params.append('id', userID);
    params.append('mID', mID);
    params.append('limit', DEFAULT_LIMIT);
    const response = yield call(instance.get, `${BASEURL}/plans/items?${params.toString()}`);
    yield put(maintenancePlanItemActions.getItemsInMaintenancePlanSuccess(response.data));
  } catch (e) {
    yield put(maintenancePlanItemActions.getItemsInMaintenancePlanFailure(e));
  }
}

export function* fetchAddItemsInPlan(action) {
  try {
    const userID = localStorage.getItem('userID');
    const { id, selectedIDList, collaborators } = action.payload;
    const response = yield call(instance.post, `${BASEURL}/plans/items`, {
      id,
      userID,
      assetIDs: selectedIDList,
      collaborators: collaborators,
    });
    yield put(maintenancePlanItemActions.addItemsInPlanSuccess(response.data));
  } catch (e) {
    yield put(maintenancePlanItemActions.addItemsInPlanFailure(e));
  }
}

export function* removeItemsFromMaintenancePlan(action) {
  try {
    const userID = localStorage.getItem('userID');
    const { id, selectedIDList } = action.payload;
    yield call(instance.post, `${BASEURL}/plan/remove/items`, {
      id,
      userID,
      assetIDs: selectedIDList,
    });
    yield put(maintenancePlanItemActions.removeItemsFromMaintenancePlanSuccess(selectedIDList));
  } catch (e) {
    yield put(maintenancePlanItemActions.removeItemsFromMaintenancePlanFailure(e));
  }
}

export function* uploadImage(action) {
  try {
    const { id, selectedImage } = action.payload;
    const formData = new FormData();
    formData.append('imageSrc', selectedImage);
    const response = yield call(instance.post, `${BASEURL}/${id}/uploadImage`, formData);
    yield put(maintenancePlanItemActions.uploadImageSuccess(response.data));
  } catch (e) {
    yield put(maintenancePlanItemActions.uploadImageFailure(e));
  }
}

export function* getSelectedImage(action) {
  try {
    const { id } = action.payload;
    // we need to modify the image to be of arrayBuffer type and build a blob object from it
    const response = yield call(instance.get, `${BASEURL}/${id}/fetchImage`, {
      responseType: 'arraybuffer',
    });

    const textDecoder = new TextDecoder();
    const responseText = textDecoder.decode(new Uint8Array(response.data));
    if (responseText.includes('NoSuchKey')) {
      throw new Error('NoSuchKey: Image does not exist');
    }

    const blob = new Blob([response.data], { type: response.headers['content-type'] });
    const avatarUrl = URL.createObjectURL(blob);
    yield put(maintenancePlanItemActions.getSelectedImageSuccess(avatarUrl));
  } catch (e) {
    yield put(maintenancePlanItemActions.getSelectedImageFailure(e));
  }
}

export function* watchGetSelectedMaintenancePlan() {
  yield takeLatest(`maintenancePlanItem/getSelectedMaintenancePlan`, getSelectedMaintenancePlan);
}

export function* watchGetItemsInMaintenancePlan() {
  yield takeLatest(`maintenancePlanItem/getItemsInMaintenancePlan`, getItemsInMaintenancePlan);
}

export function* watchUpdateCollaborators() {
  yield takeLatest(`maintenancePlanItem/updateCollaborators`, updateCollaborators);
}

export function* watchFetchAddItemsInPlan() {
  yield takeLatest(`maintenancePlanItem/addItemsInPlan`, fetchAddItemsInPlan);
}

export function* watchRemoveItemsFromMaintenancePlan() {
  yield takeLatest(`maintenancePlanItem/removeItemsFromMaintenancePlan`, removeItemsFromMaintenancePlan);
}

export function* watchUploadImage() {
  yield takeLatest(`maintenancePlanItem/uploadImage`, uploadImage);
}

export function* watchGetSelectedImage() {
  yield takeLatest(`maintenancePlanItem/getSelectedImage`, getSelectedImage);
}

export default [
  watchFetchAddItemsInPlan,
  watchGetItemsInMaintenancePlan,
  watchGetSelectedMaintenancePlan,
  watchRemoveItemsFromMaintenancePlan,
  watchGetSelectedImage,
  watchUpdateCollaborators,
  watchUploadImage,
];
