import axios from 'axios';
import { authActions } from './authSlice';

import { takeLatest, put, call, takeEvery } from 'redux-saga/effects';

import { REACT_APP_LOCALHOST_URL } from '@utils/Common';

const BASEURL = `${REACT_APP_LOCALHOST_URL}/api/v1`;

export function* fetchUserLogin(action) {
  try {
    const response = yield call(fetch, `${BASEURL}/signin`, {
      method: 'POST',
      body: JSON.stringify(action.payload),
      headers: {
        Accept: 'application/json',
      },
      credentials: 'include', // Equivalent to withCredentials: true
    });

    if (!response.ok) {
      throw new Error(`Failed to find user`);
    }

    const data = yield response.json();
    localStorage.setItem('userID', data);
    yield put(authActions.getUserIDSuccess(data));
  } catch (e) {
    yield put(authActions.getUserIDFailure(e));
  }
}

export function* fetchUserSignup(action) {
  try {
    const response = yield call(fetch, `${BASEURL}/signup`, {
      method: 'POST',
      body: JSON.stringify(action.payload),
      headers: {
        Accept: 'application/json',
      },
      credentials: 'include', // Equivalent to withCredentials: true
    });

    if (!response.ok) {
      throw new Error('Failed to create new user');
    }

    yield put(authActions.getSignupSuccess(response.data));
  } catch (e) {
    yield put(authActions.getSignupFailure(e));
  }
}

export function* fetchIsValidUserEmail(action) {
  try {
    const response = yield call(axios.post, `${BASEURL}/isValidEmail`, action.payload, {
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      },
    });
    yield put(authActions.isValidUserEmailSuccess(response.data));
  } catch (e) {
    yield put(authActions.isValidUserEmailFailure(e));
  }
}

export function* fetchUserLogout() {
  try {
    const response = yield call(axios.get, `${BASEURL}/logout`);
    yield put(authActions.getLogoutSuccess(response.data));
  } catch (e) {
    yield put(authActions.getLogoutFailure(e));
  }
}

export function* watchFetchUserLogin() {
  yield takeLatest('auth/getUserID', fetchUserLogin);
}
export function* watchFetchUserSignup() {
  yield takeLatest('auth/getSignup', fetchUserSignup);
}
export function* watchFetchIsValidUserEmail() {
  yield takeEvery(`auth/isValidUserEmail`, fetchIsValidUserEmail);
}
export function* watchFetchUserLogout() {
  yield takeLatest(`auth/getLogout`, fetchUserLogout);
}

export default [watchFetchUserLogin, watchFetchUserSignup, watchFetchIsValidUserEmail, watchFetchUserLogout];
