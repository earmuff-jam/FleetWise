import dayjs from 'dayjs';
import { useState } from 'react';

import { useDispatch, useSelector } from 'react-redux';

import {
  Box,
  Card,
  CardActions,
  CardContent,
  CardMedia,
  Checkbox,
  Skeleton,
  Stack,
  Tooltip,
  Typography,
} from '@mui/material';

import { useNavigate } from 'react-router-dom';

import SimpleModal from '@common/SimpleModal';
import ImagePicker from '@common/ImagePicker/ImagePicker';
import { inventoryActions } from '@features/Assets/inventorySlice';
import { capitalizeFirstLetter, EmptyComponent } from '@common/utils';

const GridComponent = ({ isLoading, data, rowSelected, handleRowSelection }) => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { inventories } = useSelector((state) => state.inventory);

  const [displayModal, setDisplayModal] = useState(false);
  const [selectedItemID, setSelectedItemID] = useState(-1);
  const handleCloseModal = () => setDisplayModal(false);

  const selectedAsset = inventories.filter((v) => v.id === selectedItemID).find(() => true);

  const handleNavigate = (id) => navigate(`/inventories/${id}/update`);

  const handleUpload = (id, imgFormData) => {
    dispatch(inventoryActions.uploadAndRefreshData({ id: id, selectedImage: imgFormData }));
    handleCloseModal();
  };

  if (isLoading) return <Skeleton height="10vh" />;
  if (data?.length <= 0) return <EmptyComponent />;

  return (
    <Box sx={{ overflow: 'auto' }}>
      <Stack
        spacing={{ xs: 2 }}
        marginBottom="1rem"
        direction="row"
        justifyContent={'center'}
        useFlexGap
        flexWrap="wrap"
      >
        {data.map((row) => {
          const isSelected = (id) => rowSelected.indexOf(id) !== -1;
          const selectedID = row.id;
          const isItemSelected = isSelected(selectedID);
          return (
            <Stack key={row.id} direction="row" useFlexGap flexWrap="wrap">
              <Card
                sx={{
                  display: 'flex',
                  flexDirection: 'column',
                }}
              >
                <Tooltip title={row.description}>
                  <CardMedia
                    sx={{ height: '10rem' }}
                    image={row.image || '/blank_canvas.png'}
                    onClick={() => {
                      setDisplayModal(true);
                      setSelectedItemID(selectedID);
                    }}
                  />
                </Tooltip>
                <CardContent sx={{ flexGrow: 1 }}>
                  <Stack sx={{ width: '12rem', flexGrow: 1, cursor: 'pointer' }} onClick={() => handleNavigate(row.id)}>
                    <Typography variant="h6" color="text.secondary">
                      {capitalizeFirstLetter(row.name)}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      {capitalizeFirstLetter(row.description)}
                    </Typography>
                  </Stack>
                </CardContent>
                <CardActions>
                  <Stack flexGrow={1}>
                    <Typography variant="caption">{row.quantity} items</Typography>
                    <Typography variant="caption">Edited around {dayjs(row.updated_at).fromNow()}</Typography>
                  </Stack>
                  <Checkbox
                    checked={isItemSelected}
                    color="primary"
                    size="small"
                    onClick={(event) => handleRowSelection(event, selectedID)}
                    inputProps={{ 'aria-labelledby': 'labelId' }}
                  />
                </CardActions>
              </Card>
            </Stack>
          );
        })}
        {displayModal && (
          <SimpleModal title={'Edit image'} handleClose={handleCloseModal} maxSize="xs">
            <ImagePicker id={selectedAsset.id} name={selectedAsset.name} handleUpload={handleUpload} disableCancel />
          </SimpleModal>
        )}
      </Stack>
    </Box>
  );
};

export default GridComponent;
