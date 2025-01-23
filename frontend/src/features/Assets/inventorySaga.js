import instance from '@utils/Instances';
import { Base64ToUint8Array } from '@common/utils';
import { REACT_APP_LOCALHOST_URL } from '@utils/Common';
import { inventoryActions } from '@features/Assets/inventorySlice';
import { call, put, takeEvery, takeLatest } from 'redux-saga/effects';

const BASEURL = `${REACT_APP_LOCALHOST_URL}/api/v1`;
const STORAGE_LOCATIONS_BASE_URL = `${REACT_APP_LOCALHOST_URL}/api/v1`;

export function* fetchAllInventoriesForUser(action) {
  try {
    const params = new URLSearchParams();
    const USER_ID = localStorage.getItem('userID');

    if (action.payload) {
      const { since } = action.payload;
      params.append('since', since);
    }

    const response = yield call(instance.get, `${BASEURL}/profile/${USER_ID}/inventories?${params.toString()}`);
    const data = Array.isArray(response.data) ? response.data : [];

    if (Array.isArray(data)) {
      data.map((item) => {
        if (item.image) {
          try {
            // decode the image url into a valid data
            const binaryData = Base64ToUint8Array(item.image);
            const blob = new Blob([binaryData], { type: 'image/png' });
            item.image = URL.createObjectURL(blob);
          } catch (error) {
            console.debug(error);
          }
        }
      });
    }

    yield put(inventoryActions.getAllInventoriesForUserSuccess(response.data || {}));
  } catch (e) {
    yield put(inventoryActions.getAllInventoriesForUserFailure(e));
  }
}

export function* fetchInvByID({ payload }) {
  try {
    const USER_ID = localStorage.getItem('userID');
    const response = yield call(instance.get, `${BASEURL}/profile/${USER_ID}/inventories/${payload}`);
    yield put(inventoryActions.getInvByIDSuccess(response.data || []));
  } catch (e) {
    yield put(inventoryActions.getInvByIDFailure(e));
  }
}

export function* fetchAddNewInventory(action) {
  try {
    const USER_ID = localStorage.getItem('userID');
    const response = yield call(instance.post, `${BASEURL}/profile/${USER_ID}/inventories`, {
      ...action.payload,
      created_by: USER_ID,
    });
    yield put(inventoryActions.addInventorySuccess(response.data));
  } catch (e) {
    yield put(inventoryActions.addInventoryFailure(e));
  }
}

export function* fetchAddBulkInventory(action) {
  try {
    const USER_ID = localStorage.getItem('userID');
    const response = yield call(instance.post, `${BASEURL}/profile/${USER_ID}/inventories/bulk`, { ...action.payload });
    yield put(inventoryActions.addBulkInventorySuccess(response.data));
  } catch (e) {
    yield put(inventoryActions.addBulkInventoryFailure(e));
  }
}

export function* fetchUpdateExistingInventoryDetails(action) {
  try {
    const USER_ID = localStorage.getItem('userID');
    const response = yield call(instance.put, `${BASEURL}/profile/${USER_ID}/inventories`, { ...action.payload });
    yield put(inventoryActions.updateInventorySuccess(response.data));
  } catch (e) {
    yield put(inventoryActions.updateInventoryFailure(e));
  }
}

export function* fetchUpdateAssetCol(action) {
  try {
    const USER_ID = localStorage.getItem('userID');
    const response = yield call(instance.put, `${BASEURL}/profile/${USER_ID}/inventories/${action.payload.assetID}`, {
      ...action.payload,
    });
    yield put(inventoryActions.updateAssetColSuccess(response.data || []));
  } catch (e) {
    yield put(inventoryActions.updateAssetColFailure(e));
  }
}

export function* fetchRemoveInventoryRows(action) {
  try {
    const USER_ID = localStorage.getItem('userID');
    const response = yield call(instance.post, `${BASEURL}/profile/${USER_ID}/inventories/prune`, {
      ...action.payload,
    });
    yield put(inventoryActions.removeInventoryRowsSuccess(response.data));
  } catch (e) {
    yield put(inventoryActions.removeInventoryRowsFailure(e));
  }
}

export function* fetchStorageLocations() {
  try {
    const response = yield call(instance.get, `${STORAGE_LOCATIONS_BASE_URL}/locations`);
    yield put(inventoryActions.getStorageLocationsSuccess(response.data || []));
  } catch (e) {
    yield put(inventoryActions.getStorageLocationsFailure(e));
  }
}

export function* uploadImage(action) {
  try {
    const { id, selectedImage } = action.payload;
    const formData = new FormData();
    formData.append('imageSrc', selectedImage);
    const response = yield call(instance.post, `${BASEURL}/${id}/uploadImage`, formData);
    yield put(inventoryActions.uploadImageSuccess(response.data));
  } catch (e) {
    yield put(inventoryActions.uploadImageFailure(e));
  }
}

export function* uploadAndRefreshDataSuccess(action) {
  /**
   * Fixes https://github.com/earmuff-jam/FleetWise/issues/305.
   * Refetch the data once the image is updated in the system. Re-use the same function for selected asset details.
   * If the user is viewing selected asset and refreshes the page, the user should see the new image as well.
   */
  try {
    const { id, selectedImage } = action.payload;
    const formData = new FormData();
    formData.append('imageSrc', selectedImage);

    yield call(instance.post, `${BASEURL}/${id}/uploadImage`, formData);

    const params = new URLSearchParams();
    const USER_ID = localStorage.getItem('userID');
    const response = yield call(instance.get, `${BASEURL}/profile/${USER_ID}/inventories?${params.toString()}`);
    const data = Array.isArray(response.data) ? response.data : [];

    if (Array.isArray(data)) {
      data.map((item) => {
        if (item.image) {
          try {
            // decode the image url into a valid data
            const binaryData = Base64ToUint8Array(item.image);
            const blob = new Blob([binaryData], { type: 'image/png' });
            item.image = URL.createObjectURL(blob);
          } catch (error) {
            console.debug(error);
          }
        }
      });
    }

    yield put(inventoryActions.getAllInventoriesForUserSuccess(response.data || {}));
    yield put(inventoryActions.getSelectedImage({ id: id }));
    yield put(inventoryActions.uploadImageSuccess(response.data));
  } catch (e) {
    yield put(inventoryActions.uploadImageFailure(e));
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
    yield put(inventoryActions.getSelectedImageSuccess(avatarUrl));
  } catch (e) {
    yield put(inventoryActions.getSelectedImageFailure(e));
  }
}

/********************************
 * WATCHER FUNCTIONS BELOW HERE
 ********************************/

export function* watchFetchAllInventoriesForUser() {
  yield takeLatest(`inventory/getAllInventoriesForUser`, fetchAllInventoriesForUser);
}

export function* watchFetchInvByID() {
  yield takeLatest(`inventory/getInvByID`, fetchInvByID);
}

export function* watchFetchAddNewInventory() {
  yield takeLatest(`inventory/addInventory`, fetchAddNewInventory);
}

export function* watchFetchAddBulkInventory() {
  yield takeLatest(`inventory/addBulkInventory`, fetchAddBulkInventory);
}

export function* watchUpdateExistingInventoryDetails() {
  yield takeLatest(`inventory/updateInventory`, fetchUpdateExistingInventoryDetails);
}

export function* watchFetchRemoveInventoryRows() {
  yield takeLatest(`inventory/removeInventoryRows`, fetchRemoveInventoryRows);
}

export function* watchFetchAllStorageLocations() {
  yield takeLatest(`inventory/getStorageLocations`, fetchStorageLocations);
}

export function* watchUpdateAssetCol() {
  yield takeLatest(`inventory/updateAssetCol`, fetchUpdateAssetCol);
}

export function* watchUploadImage() {
  yield takeEvery(`inventory/uploadImage`, uploadImage);
}

export function* watchUploadAndRefreshDataSuccess() {
  yield takeEvery(`inventory/uploadAndRefreshData`, uploadAndRefreshDataSuccess);
}

export function* watchGetSelectedImage() {
  yield takeLatest(`inventory/getSelectedImage`, getSelectedImage);
}

export default [
  watchFetchInvByID,
  watchUpdateAssetCol,
  watchFetchAddNewInventory,
  watchUploadImage,
  watchGetSelectedImage,
  watchFetchAddBulkInventory,
  watchFetchRemoveInventoryRows,
  watchFetchAllStorageLocations,
  watchFetchAllInventoriesForUser,
  watchUploadAndRefreshDataSuccess,
  watchUpdateExistingInventoryDetails,
];
