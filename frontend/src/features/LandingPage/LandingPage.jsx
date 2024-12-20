import { useState } from 'react';

import { Container, Stack } from '@mui/material';

import SimpleModal from '@common/SimpleModal';
import FloatingBarChart from '@common/Chart/FloatingBarChart';
import Login from '@features/LandingPage/Authentication/Login/Login';

import { COLORS, SAMPLE_DATA } from '@features/LandingPage/constants';
import HeroContent from '@features/LandingPage/HeroContent/HeroContent';
import Signup from '@features/LandingPage/Authentication/Signup/Signup';
import StyledAppBar from '@features/LandingPage/StyledAppBar/StyledAppBar';

export default function LandingPage() {
  const [formType, setFormType] = useState('');
  const [displayModal, setDisplayModal] = useState(false);

  const openDisplayModal = () => setDisplayModal(true);
  const closeDisplayModal = () => setDisplayModal(false);

  const openSignupModal = () => {
    openDisplayModal();
    setFormType('signup');
  };

  const openLoginModal = () => {
    openDisplayModal();
    setFormType('signin');
  };

  const formattedData = SAMPLE_DATA.map((v, index) => ({
    label: v.name,
    start: index === 0 ? 0 : SAMPLE_DATA[index - 1].price,
    end: v.price,
    color: COLORS[index % COLORS.length],
  }));

  return (
    <>
      <StyledAppBar />
      <Container maxWidth="md">
        <Stack spacing={2}>
          <HeroContent openSignupModal={openSignupModal} openLoginModal={openLoginModal} />
          <FloatingBarChart
            data={formattedData}
            backgroundColor={formattedData.map((d) => d.color)}
            borderColor={formattedData.map((d) => d.color)}
          />
        </Stack>
        {displayModal && (
          <SimpleModal
            title={formType === 'signin' ? 'Sign In' : 'Sign Up'}
            subtitle={
              formType === 'signin'
                ? 'Login and manage your account.'
                : 'Create an account to keep track of all your inventories.'
            }
            handleClose={closeDisplayModal}
            maxSize="xs"
          >
            {formType === 'signup' ? <Signup handleClose={closeDisplayModal} /> : <Login />}
          </SimpleModal>
        )}
      </Container>
    </>
  );
}
