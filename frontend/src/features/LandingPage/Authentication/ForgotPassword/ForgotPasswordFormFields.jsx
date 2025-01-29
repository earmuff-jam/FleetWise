import { InputAdornment, Stack, TextField, Typography } from '@mui/material';

export default function ForgotPasswordFormFields({ formFields, handleInput }) {
  return (
    <Stack spacing={1}>
      <Typography variant="subtitle2" color="text.secondary">
        {formFields['email'].label} {formFields['email'].required && '*'}
      </Typography>
      <TextField
        id={formFields['email'].name}
        name={formFields['email'].name}
        size={formFields['email'].size}
        value={formFields['email'].value}
        type={formFields['email'].type}
        variant={formFields['email'].variant}
        placeholder={formFields['email'].placeholder}
        onChange={handleInput}
        required={formFields['email'].required}
        fullWidth={formFields['email'].fullWidth}
        error={!!formFields['email'].errorMsg}
        helperText={formFields['email'].errorMsg}
        InputProps={{
          startAdornment: <InputAdornment position="start">{formFields['email'].icon}</InputAdornment>,
        }}
      />
    </Stack>
  );
}
