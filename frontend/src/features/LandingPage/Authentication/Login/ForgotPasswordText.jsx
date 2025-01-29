import { Stack, Typography } from '@mui/material';

export default function ForgotPasswordText({ handleForgotPassword }) {
  return (
    <Stack direction="row" alignItems="center" spacing={0.4}>
      <Typography variant="caption">Forgot your password? </Typography>
      <Typography
        variant="caption"
        sx={{ cursor: 'pointer', '&:hover': { color: 'info.main' } }}
        onClick={handleForgotPassword}
      >
        Reset it here.
      </Typography>
    </Stack>
  );
}
