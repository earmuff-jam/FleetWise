import { produce } from 'immer';
import { useState } from 'react';
import { useDispatch } from 'react-redux';
import { Button, Stack } from '@mui/material';
import { authActions } from '@features/LandingPage/authSlice';
import { LOGIN_FORM_FIELDS } from '@features/LandingPage/constants';
import LoginFormFields from '@features/LandingPage/Authentication/Login/LoginFormFields';
import ForgotPasswordText from '@features/LandingPage/Authentication/Login/ForgotPasswordText';

const Login = ({ handleForgotPassword }) => {
  const dispatch = useDispatch();
  const [formFields, setFormFields] = useState(LOGIN_FORM_FIELDS);

  const handleInput = (event) => {
    const { name, value } = event.target;
    setFormFields(
      produce(formFields, (draft) => {
        draft[name].value = value;
        draft[name].errorMsg = '';

        for (const validator of draft[name].validators) {
          if (validator.validate(value)) {
            draft[name].errorMsg = validator.message;
            break;
          }
        }
      })
    );
  };

  const validate = (formFields) => {
    const containsErr = Object.values(formFields).reduce((acc, el) => {
      if (el.errorMsg) {
        return true;
      }
      return acc;
    }, false);

    const requiredFormFields = Object.values(formFields).filter((v) => v.required);
    return containsErr || requiredFormFields.some((el) => el.value.trim() === '');
  };

  const submit = (e) => {
    e.preventDefault();

    if (validate(formFields)) {
      return;
    } else {
      const formattedData = Object.values(formFields).reduce((acc, el) => {
        if (el.value) {
          acc[el.name] = el.value;
        }
        return acc;
      }, {});
      dispatch(authActions.getUserID(formattedData));
    }
  };

  return (
    <Stack spacing={1}>
      <LoginFormFields formFields={formFields} handleInput={handleInput} submit={submit} />
      <ForgotPasswordText handleForgotPassword={handleForgotPassword} />
      <Button
        variant="text"
        disabled={validate(formFields)}
        disableRipple={true}
        disableFocusRipple={true}
        onClick={submit}
      >
        Login
      </Button>
    </Stack>
  );
};

export default Login;
