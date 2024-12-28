import dayjs from 'dayjs';

import { AddPhotoAlternateRounded, ShareRounded } from '@mui/icons-material';
import { Badge, Button, CardActions, IconButton, Stack, Typography } from '@mui/material';

export default function DetailsCardItemActions({
  selectedItem,
  handleOpenModal,
  setEditImgMode,
  imageBtnDataTour,
  shareBtnDataTour,
}) {
  const isShared = selectedItem?.sharable_groups?.length > 1 || false;

  return (
    <CardActions
      sx={{ display: 'flex', flexDirection: 'row', justifyContent: 'space-between', alignContent: 'center' }}
    >
      <Typography variant="caption">Last updated {dayjs(selectedItem?.updated_at).fromNow()}</Typography>
      <Stack direction="row" alignItems="center">
        <Badge
          badgeContent={isShared ? selectedItem?.sharable_groups.length - 1 : 0} // account for creator in sharable_groups
          color="secondary"
          max={10}
          anchorOrigin={{
            vertical: 'top',
            horizontal: 'left',
          }}
        >
          <Button size="small" endIcon={<ShareRounded />} onClick={handleOpenModal} data-tour={shareBtnDataTour}>
            Share
          </Button>
        </Badge>
        <IconButton onClick={() => setEditImgMode(true)} data-tour={imageBtnDataTour}>
          <AddPhotoAlternateRounded color="primary" />
        </IconButton>
      </Stack>
    </CardActions>
  );
}
