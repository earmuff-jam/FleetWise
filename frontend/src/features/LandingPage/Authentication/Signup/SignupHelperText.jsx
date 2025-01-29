import { CircularProgress, Stack, Typography } from '@mui/material';
import { CheckRounded, CloseRounded } from '@mui/icons-material';

export default function SignupHelperText({ isEmailUnique = false, loading = false }) {
  if (loading) {
    return <CircularProgress size="1.2rem" />;
  }

  return isEmailUnique ? (
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
