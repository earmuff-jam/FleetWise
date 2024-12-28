import { useDispatch, useSelector } from 'react-redux';

import { FavoriteRounded } from '@mui/icons-material';
import { CardContent, IconButton, Stack, Typography } from '@mui/material';
import { profileActions } from '@features/Profile/profileSlice';

export default function DetailsCardItemContent({ selectedItem, categoryMode, favBtnDataTour }) {
  const dispatch = useDispatch();
  const { favItems = [] } = useSelector((state) => state.profile);

  const isFavourite = favItems.some(
    (v) => v.category_id === selectedItem.id || v.maintenance_plan_id === selectedItem.id
  );

  const handleFavItem = (_, selectedID, isFavourite) => {
    let draftFavItem = {};
    if (categoryMode) {
      draftFavItem = { category_id: selectedID };
    } else {
      draftFavItem = { maintenance_plan_id: selectedID };
    }

    if (isFavourite) {
      // toggle fav off if exists
      const currentItems = favItems?.filter(
        (v) => v.category_id === selectedID || v.maintenance_plan_id === selectedID
      );
      const currentItem = currentItems.find(() => true);
      dispatch(profileActions.removeFavItem(currentItem?.id));
    } else {
      dispatch(profileActions.saveFavItem(draftFavItem));
    }
  };

  return (
    <CardContent>
      <Stack direction="row" alignItems="flex-start">
        <IconButton
          size="small"
          onClick={(ev) => handleFavItem(ev, selectedItem.id, isFavourite)}
          data-tour={favBtnDataTour}
        >
          <FavoriteRounded fontSize="small" sx={{ color: isFavourite ? selectedItem.color : 'secondary.main' }} />
        </IconButton>
        <Typography gutterBottom variant="h5">
          {selectedItem.name}
        </Typography>
      </Stack>
      <Typography variant="body2" color="text.secondary">
        {selectedItem.description}
      </Typography>
    </CardContent>
  );
}
