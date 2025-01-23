import { useState, useEffect } from 'react';
import { Alert, Collapse, IconButton } from '@mui/material';
import EmergencyShare from '@mui/icons-material/EmergencyShare';

export default function Banner({ isVerified = 'false', revalidateEmail }) {
  const [timer, setTimer] = useState(0);

  const disabledEmailValidation = localStorage.getItem('disableEmailValidation');

  useEffect(() => {
    let countdown;
    if (timer > 0) {
      countdown = setInterval(() => {
        setTimer((prev) => prev - 1);
      }, 1000);
    }

    return () => {
      clearInterval(countdown);
    };
  }, [timer]);

  const handleClick = () => {
    setTimer(60);
    localStorage.setItem('disableEmailValidation', true);
    revalidateEmail();
  };

  return (
    <Collapse in={isVerified === 'false'}>
      <Alert
        severity="warning"
        sx={{ margin: '1rem 0rem' }}
        action={
          <IconButton
            aria-label="close"
            color="inherit"
            size="small"
            onClick={handleClick}
            disabled={timer > 0 || disabledEmailValidation}
          >
            <EmergencyShare fontSize="inherit" />
          </IconButton>
        }
      >
        Please verify your email address
        {timer > 0 && <span style={{ marginLeft: '1rem' }}>Resend available in {timer} seconds</span>}
      </Alert>
    </Collapse>
  );
}
