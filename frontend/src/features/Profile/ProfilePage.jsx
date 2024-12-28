import { useState } from 'react';
import { useTheme } from '@emotion/react';
import { useSelector } from 'react-redux';
import SimpleModal from '@common/SimpleModal';
import { Paper, Stack, useMediaQuery } from '@mui/material';

import UserStatus from '@features/Profile/UserStatus/UserStatus';
import ProfileForm from '@features/Profile/ProfileForm/ProfileForm';

import UserDetails from '@features/Profile/UserDetails/UserDetails';
import UserDemographics from '@features/Profile/UserDemographics/UserDemographics';
import AppearanceSettings from '@features/Profile/AppearanceSettings/AppearanceSettings';

const ProfilePage = () => {
  const theme = useTheme();
  const onlySmallScreen = useMediaQuery(theme.breakpoints.down('sm'));
  const { profileDetails: data = {}, profileStats = {}, loading } = useSelector((state) => state.profile);

  const [editMode, setEditMode] = useState(false);

  return (
    <Stack spacing={2} data-tour="profile-0">
      <Stack direction={onlySmallScreen ? 'column' : 'row'} spacing={1}>
        <UserDemographics data={data} handleEditMode={() => setEditMode(!editMode)} />
        <UserDetails data={data} />
      </Stack>
      <Paper sx={{ padding: '1rem' }} data-tour="profile-5">
        <UserStatus profileStats={profileStats} onlySmallScreen={onlySmallScreen} />
      </Paper>
      <Paper sx={{ padding: '1rem' }} data-tour="profile-6">
        <AppearanceSettings loading={loading} profileDetails={data} />
      </Paper>
      {editMode && (
        <SimpleModal
          open={editMode}
          handleClose={() => setEditMode(false)}
          maxSize={'xs'}
          title="Edit profile details"
          subtitle="Edit general details about yourself so others can notice you."
        >
          <ProfileForm />
        </SimpleModal>
      )}
    </Stack>
  );
};

export default ProfilePage;
