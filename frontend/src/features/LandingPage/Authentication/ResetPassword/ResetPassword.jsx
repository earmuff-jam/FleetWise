import { Box, Button } from '@mui/material';

import RowHeader from '@common/RowHeader';
import ResetPasswordFormFields from '@features/LandingPage/Authentication/ResetPassword/ResetPasswordFormFields';

export default function ResetPassword() {
  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center',
        minHeight: '100vh',
        gap: 2,
      }}
    >
      <RowHeader title={'Reset your password'} caption={'Enter the OTP provided to your email address'} />
      <ResetPasswordFormFields />
      <Button variant="outlined">Submit</Button>
    </Box>
  );
}
