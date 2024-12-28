import dayjs from 'dayjs';

import relativeTime from 'dayjs/plugin/relativeTime';
import { Box, Divider, Paper, Stack, Typography } from '@mui/material';
import { GridViewRounded, OnlinePredictionRounded, UpdateRounded, ViewListRounded } from '@mui/icons-material';

dayjs.extend(relativeTime);

const BIO_PLACEHOLDER = 'Add short bio about yourself to let others know details about you.';

export default function UserDetails({ data = {} }) {
  return (
    <Paper
      sx={{
        padding: '1rem',
        display: 'flex',
        flexDirection: 'column',
        flexGrow: 1,
      }}
      data-tour="profile-3"
    >
      <Typography variant="h5" color="text.secondary">
        User details
      </Typography>
      <Divider sx={{ marginTop: '1rem', marginBottom: '1rem' }} />
      <Typography variant="caption" color="text.secondary">
        {data.about_me || BIO_PLACEHOLDER}
      </Typography>
      <Box sx={{ flexGrow: 1 }} />
      <Stack spacing={1} data-tour="profile-4">
        <Stack direction={'row'} alignItems={'center'} spacing={1}>
          <OnlinePredictionRounded color="primary" />
          <Typography variant="subtitle2" color="text.secondary">
            Online
          </Typography>
        </Stack>
        <Stack direction={'row'} alignItems={'center'} spacing={1}>
          {data?.grid_view ? <GridViewRounded color="primary" /> : <ViewListRounded color="primary" />}
          <Typography variant="subtitle2" color="text.secondary">
            {data?.grid_view ? 'Viewing Assets as Grid' : 'Viewing Assets as List'}
          </Typography>
        </Stack>
        <Stack direction={'row'} alignItems={'center'} spacing={1}>
          <UpdateRounded color="primary" />
          <Typography variant="subtitle2" color="text.secondary">
            Last updated at : {dayjs(data.updated_at).fromNow()}
          </Typography>
        </Stack>
      </Stack>
    </Paper>
  );
}
