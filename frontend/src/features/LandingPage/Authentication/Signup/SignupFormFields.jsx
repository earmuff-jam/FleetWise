import { produce } from 'immer';
import { useDispatch, useSelector } from 'react-redux';

import { CheckRounded, CloseRounded } from '@mui/icons-material';

import { Stack, TextField, InputAdornment, Typography, CircularProgress } from '@mui/material';

import { authActions } from '@features/LandingPage/authSlice';

export default function SignupFormFields({ formFields, setFormFields, submit }) {
  const dispatch = useDispatch();
  const { isValidUserEmail, loading } = useSelector((state) => state.auth);

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

  const validUserEmail = (isValidUserEmail, loading) => {
    if (loading) {
      return <CircularProgress size="1.2rem" />;
    } else {
      return isValidUserEmail ? (
        <Stack direction="row" alignItems="center" component="span" spacing={0.2}>
          <CheckRounded color="success" fontSize="small" />
          <Typography variant="caption">Unique value for email</Typography>
        </Stack>
      ) : (
        <Stack direction="row" alignItems="center" component="span" spacing={0.2}>
          <CloseRounded color="error" fontSize="small" />
          <Typography variant="caption">Existing value for email</Typography>
        </Stack>
      );
    }
  };

  return (
    <Stack spacing={2}>
      <TextField
        id={formFields['username'].name}
        name={formFields['username'].name}
        label={formFields['username'].label}
        value={formFields['username'].value}
        size={formFields['username'].size}
        type={formFields['username'].type}
        variant={formFields['username'].variant}
        autoComplete={formFields['username'].autocomplete}
        placeholder={formFields['username'].placeholder}
        onChange={handleInput}
        required={formFields['username'].required}
        fullWidth={formFields['username'].fullWidth}
        error={!!formFields['username'].errorMsg}
        helperText={formFields['username'].errorMsg}
        onKeyDown={(e) => {
          if (e.code === 'Enter') {
            submit(e);
          }
        }}
        InputProps={{
          startAdornment: <InputAdornment position="start">{formFields['username'].icon}</InputAdornment>,
        }}
      />
      <TextField
        id={formFields['email'].name}
        name={formFields['email'].name}
        label={formFields['email'].label}
        value={formFields['email'].value}
        size={formFields['email'].size}
        type={formFields['email'].type}
        variant={formFields['email'].variant}
        autoComplete={formFields['email'].autocomplete}
        placeholder={formFields['email'].placeholder}
        onChange={handleInput}
        required={formFields['email'].required}
        fullWidth={formFields['email'].fullWidth}
        error={!!formFields['email'].errorMsg}
        helperText={formFields['email'].errorMsg || validUserEmail(isValidUserEmail, loading)}
        onKeyDown={(e) => {
          if (e.code === 'Enter') {
            submit(e);
          }
        }}
        onBlur={() =>
          formFields.email.value.length > 3 && dispatch(authActions.isValidUserEmail({ email: formFields.email.value }))
        }
        InputProps={{
          startAdornment: <InputAdornment position="start">{formFields['email'].icon}</InputAdornment>,
        }}
      />
      <TextField
        id={formFields['password'].name}
        name={formFields['password'].name}
        label={formFields['password'].label}
        value={formFields['password'].value}
        size={formFields['password'].size}
        type={formFields['password'].type}
        variant={formFields['password'].variant}
        autoComplete={formFields['password'].autocomplete}
        placeholder={formFields['password'].placeholder}
        onChange={handleInput}
        required={formFields['password'].required}
        fullWidth={formFields['password'].fullWidth}
        error={!!formFields['password'].errorMsg}
        helperText={formFields['password'].errorMsg}
        onKeyDown={(e) => {
          if (e.code === 'Enter') {
            submit(e);
          }
        }}
        InputProps={{
          startAdornment: <InputAdornment position="start">{formFields['password'].icon}</InputAdornment>,
        }}
      />
      <TextField
        id={formFields['birthday'].name}
        name={formFields['birthday'].name}
        label={formFields['birthday'].label}
        value={formFields['birthday'].value}
        type={formFields['birthday'].type}
        size={formFields['birthday'].size}
        variant={formFields['birthday'].variant}
        autoComplete={formFields['birthday'].autocomplete}
        placeholder={formFields['birthday'].placeholder}
        onChange={handleInput}
        required={formFields['birthday'].required}
        fullWidth={formFields['birthday'].fullWidth}
        error={!!formFields['birthday'].errorMsg}
        helperText={formFields['birthday'].errorMsg}
        onKeyDown={(e) => {
          if (e.code === 'Enter') {
            submit(e);
          }
        }}
        InputProps={{
          startAdornment: <InputAdornment position="start">{formFields['birthday'].icon}</InputAdornment>,
        }}
      />
    </Stack>
  );
}
