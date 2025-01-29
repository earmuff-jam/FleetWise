import { useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { enqueueSnackbar } from 'notistack';
import { Button, Stack } from '@mui/material';
import { authActions } from '@features/LandingPage/authSlice';
import { SIGN_UP_FORM_FIELDS } from '@features/LandingPage/constants';
import SignupFormFields from '@features/LandingPage/Authentication/Signup/SignupFormFields';
import SingupCheckboxField from '@features/LandingPage/Authentication/Signup/SignupCheckboxField';
import SignupTermsConditions from '@features/LandingPage/Authentication/Signup/SignupTermsConditions';

export default function Signup({ handleClose }) {
  const dispatch = useDispatch();
  const { isValidUserEmail, loading } = useSelector((state) => state.auth);

  const [isChecked, setIsChecked] = useState(false);
  const [formFields, setFormFields] = useState(SIGN_UP_FORM_FIELDS);

  const hasError = (formFields = []) => {
    const containsErr = Object.values(formFields).reduce((acc, el) => {
      if (el.errorMsg) {
        return true;
      }
      return acc;
    }, false);
    const requiredFormFields = Object.values(formFields).filter((v) => v.required);
    const isRequiredFieldsEmpty = requiredFormFields.some((el) => el.value.trim() === '');

    return containsErr || isRequiredFieldsEmpty;
  };

  const submit = (e) => {
    e.preventDefault();
    if (hasError(formFields) || !isChecked) {
      return;
    } else {
      const formattedData = Object.values(formFields).reduce((acc, el) => {
        if (el.value) {
          acc[el.name] = el.value;
        }
        return acc;
      }, {});
      dispatch(authActions.getSignup(formattedData));
      enqueueSnackbar('Successfully created new user. Please login.', {
        variant: 'success',
      });
      handleClose(false);
    }
  };

  return (
    <Stack>
      <SignupFormFields
        formFields={formFields}
        setFormFields={setFormFields}
        isValidUserEmail={isValidUserEmail}
        loading={loading}
        submit={submit}
      />
      <SingupCheckboxField isChecked={isChecked} setIsChecked={setIsChecked} />
      <SignupTermsConditions />
      <Button variant="text" onClick={submit} disabled={hasError(formFields) || !isChecked || !isValidUserEmail}>
        Register
      </Button>
    </Stack>
  );
}
