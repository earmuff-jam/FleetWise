import { useState } from 'react';
import { enqueueSnackbar } from 'notistack';

import dayjs from 'dayjs';
import { produce } from 'immer';
import { useDispatch } from 'react-redux';
import { useNavigate } from 'react-router-dom';

import relativeTime from 'dayjs/plugin/relativeTime';
import { Card, CardMedia, Paper, Stack } from '@mui/material';

import SimpleModal from '@common/SimpleModal';
import SharableGroups from '@common/SharableGroups';
import ImagePicker from '@common/ImagePicker/ImagePicker';
import LocationPicker from '@common/Location/LocationPicker';

import DetailsCardItemContent from '@common/ItemCard/ItemContent/DetailsCardItemContent';
import DetailsCardItemActions from '@common/ItemCard/ItemContent/DetailsCardItemActions';
import { categoryItemDetailsActions } from '@features/CategoryItemDetails/categoryItemDetailsSlice';
import { maintenancePlanItemActions } from '@features/MaintenancePlanItemDetails/maintenancePlanItemSlice';

dayjs.extend(relativeTime);

export default function DetailsCard({
  selectedItem,
  selectedImage,
  favBtnDataTour,
  imageBtnDataTour,
  shareBtnDataTour,
  categoryMode = false,
}) {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const userID = localStorage.getItem('userID');

  const [editImgMode, setEditImgMode] = useState(null);
  const [openModal, setOpenModal] = useState(false);

  const handleOpenModal = () => setOpenModal(true);
  const handleCloseModal = () => setOpenModal(false);

  const handleUpload = (id, selectedImage) => {
    if (categoryMode) {
      dispatch(categoryItemDetailsActions.uploadImage({ id: id, selectedImage: selectedImage }));
    } else {
      dispatch(maintenancePlanItemActions.uploadImage({ id: id, selectedImage: selectedImage }));
    }
    enqueueSnackbar('New image upload successful.', {
      variant: 'success',
    });
    setEditImgMode(false);
  };

  const updateCollaborators = (sharableGroups) => {
    const newMembers = sharableGroups.map((v) => v.value);
    const draftSelectionDetails = produce(selectedItem, (draft) => {
      draft.updated_by = userID;
      draft.sharable_groups = newMembers;
      if (categoryMode) {
        draft.status = draft.status_name;
      } else {
        draft.maintenance_status = draft.maintenance_status_name;
      }
    });
    if (categoryMode) {
      dispatch(categoryItemDetailsActions.updateCollaborators(draftSelectionDetails));
      enqueueSnackbar('Updated collaborators for selected category.', {
        variant: 'success',
      });
      if (!newMembers.includes(userID)) {
        navigate('/');
      }
    } else {
      dispatch(maintenancePlanItemActions.updateCollaborators(draftSelectionDetails));
      enqueueSnackbar('Updated collaborators for selected maintenance plan.', {
        variant: 'success',
      });
      if (!newMembers.includes(userID)) {
        navigate('/');
      }
    }
    handleCloseModal();
  };

  return (
    <>
      <Stack sx={{ flexDirection: { xs: 'column', sm: 'row' }, gap: 2 }}>
        <Card sx={{ flexGrow: 1 }}>
          <CardMedia sx={{ height: '10rem' }} image={selectedImage || '/blank_canvas.png'} />
          <DetailsCardItemContent
            selectedItem={selectedItem}
            categoryMode={categoryMode}
            favBtnDataTour={favBtnDataTour}
          />
          <DetailsCardItemActions
            selectedItem={selectedItem}
            handleOpenModal={handleOpenModal}
            setEditImgMode={setEditImgMode}
            imageBtnDataTour={imageBtnDataTour}
            shareBtnDataTour={shareBtnDataTour}
          />
        </Card>
        {selectedItem?.location?.lat ? (
          <Paper elevation={2} sx={{ width: { xs: '100%', sm: '20rem' }, height: 'inherit', flexGrow: 1 }}>
            <LocationPicker location={selectedItem?.location} height={'100%'} />
          </Paper>
        ) : null}
      </Stack>
      {openModal && (
        <SimpleModal
          title="Add sharable groups"
          subtitle="Assign collaborators."
          handleClose={handleCloseModal}
          maxSize="sm"
        >
          <SharableGroups
            handleSubmit={updateCollaborators}
            existingGroups={selectedItem?.sharable_groups || []}
            creator={selectedItem?.created_by}
          />
        </SimpleModal>
      )}
      {editImgMode && (
        <SimpleModal
          title="Assign image"
          subtitle="Assign image to the selected item."
          handleClose={() => setEditImgMode(false)}
          maxSize="xs"
        >
          <ImagePicker id={selectedItem.id} name={selectedItem.name} handleUpload={handleUpload} disableCancel />
        </SimpleModal>
      )}
    </>
  );
}
