import { Stack, TextField, InputAdornment, Typography } from '@mui/material';

export default function LoginFormFields({ formFields, handleInput, submit }) {
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
        autoComplete={formFields['email'].autocomplete}
        placeholder={formFields['email'].placeholder}
        onChange={handleInput}
        required={formFields['email'].required}
        fullWidth={formFields['email'].fullWidth}
        error={!!formFields['email'].errorMsg}
        helperText={formFields['email'].errorMsg}
        onKeyDown={(e) => {
          if (e.code === 'Enter') {
            submit(e);
          }
        }}
        InputProps={{
          startAdornment: <InputAdornment position="start">{formFields['email'].icon}</InputAdornment>,
        }}
      />

      <Typography variant="subtitle2" color="text.secondary">
        {formFields['password'].label} {formFields['password'].required && '*'}
      </Typography>
      <TextField
        id={formFields['password'].name}
        name={formFields['password'].name}
        size={formFields['password'].size}
        value={formFields['password'].value}
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
    </Stack>
  );
}
