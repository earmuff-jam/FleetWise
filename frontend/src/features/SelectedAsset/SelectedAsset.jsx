import { useEffect, useState } from 'react';

import dayjs from 'dayjs';
import { enqueueSnackbar } from 'notistack';
import { useDispatch, useSelector } from 'react-redux';
import { useParams } from 'react-router-dom';
import RowHeader from '@common/RowHeader';
import SimpleModal from '@common/SimpleModal';
import relativeTime from 'dayjs/plugin/relativeTime';
import ImagePicker from '@common/ImagePicker/ImagePicker';
import { Button, Divider, Stack, Typography } from '@mui/material';
import { BLANK_INVENTORY_FORM } from '@features/Assets/constants';
import { inventoryActions } from '@features/Assets/inventorySlice';
import { AddPhotoAlternateRounded, CheckRounded } from '@mui/icons-material';
import SelectedAssetFormFields from '@features/SelectedAsset/SelectedAssetFormFields';
import SelectedAssetWeightDimension from '@features/SelectedAsset/SelectedAssetWeightDimension';
import SelectedAssetMoreInformation from '@features/SelectedAsset/SelectedAssetMoreInformation';

dayjs.extend(relativeTime);

export default function SelectedAsset() {
  const { id } = useParams();
  const dispatch = useDispatch();

  const {
    loading: storageLocationsLoading,
    storageLocations,
    inventory,
    selectedImage,
    loading,
  } = useSelector((state) => state.inventory);

  const [editImgMode, setEditImgMode] = useState(false);
  const [openReturnNote, setOpenReturnNote] = useState(false);
  const [formFieldChange, setFormFieldChange] = useState(false);

  const [color, setColor] = useState('#f7f7f7');
  const [returnDateTime, setReturnDateTime] = useState(null);
  const [storageLocation, setStorageLocation] = useState({});
  const [formData, setFormData] = useState({ ...BLANK_INVENTORY_FORM });

  const handleInputChange = (event) => {
    const { id, value } = event.target;
    const updatedFormData = { ...formData };
    let errorMsg = '';

    for (const validator of updatedFormData[id].validators) {
      if (validator.validate(value)) {
        errorMsg = validator.message;
        break;
      }
    }

    updatedFormData[id] = {
      ...updatedFormData[id],
      value,
      errorMsg,
    };
    setFormFieldChange(true);
    setFormData(updatedFormData);
  };

  const handleColorChange = (newValue) => {
    setFormFieldChange(true);
    setColor(newValue);
  };

  const handleCheckbox = (name, value) => {
    // close the return note subsection if the parent is disabled
    if (name === 'is_returnable' && openReturnNote) {
      setOpenReturnNote(false);
    }

    setFormFieldChange(true);
    setFormData((prevFormData) => ({
      ...prevFormData,
      [name]: { ...prevFormData[name], value },
    }));
  };

  const isFormDisabled = () => {
    const containsErr = Object.values(formData).reduce((acc, el) => {
      if (el?.errorMsg) {
        return true;
      }
      return acc;
    }, false);

    const requiredFormFields = Object.values(formData).filter((v) => v?.isRequired);
    const isRequiredFieldsEmpty = requiredFormFields
      .filter((el) => el.type === 'text')
      .some((el) => el.value.trim() === '');

    return containsErr || isRequiredFieldsEmpty || storageLocation === null || Object.keys(storageLocation).length <= 0;
  };

  const handleSubmit = (event) => {
    event.preventDefault();

    if (isFormDisabled()) {
      enqueueSnackbar('Unable to update inventory details.', {
        variant: 'error',
      });
      return;
    }

    const formattedData = Object.values(formData).reduce((acc, el) => {
      if (el.type === 'number') {
        acc[el.id] = Number(el.value);
      } else if (el.value) {
        acc[el.id] = el.value;
      }
      return acc;
    }, {});

    const draftRequest = {
      id: id, // bring id from the params
      ...formattedData,
      return_datetime: returnDateTime !== null ? returnDateTime.toISOString() : null,
      location: storageLocation.location,
      color: color,
    };

    setFormFieldChange(false);
    dispatch(inventoryActions.updateInventory(draftRequest));
    enqueueSnackbar('Successfully updated selected asset.', {
      variant: 'success',
    });
  };

  const handleUpload = (id, imgFormData) => {
    dispatch(inventoryActions.uploadAndRefreshData({ id: id, selectedImage: imgFormData }));
    setEditImgMode(false);
  };

  useEffect(() => {
    if (id.length > 0) {
      dispatch(inventoryActions.getInvByID(id));
    }
  }, [id]);

  useEffect(() => {
    if (!loading && !selectedImage) {
      dispatch(inventoryActions.getSelectedImage({ id }));
    }
  }, [loading]);

  useEffect(() => {
    if (!loading || !storageLocationsLoading) {
      const selectedAsset = { ...BLANK_INVENTORY_FORM };
      selectedAsset.name.value = inventory.name || '';
      selectedAsset.description.value = inventory.description || '';
      selectedAsset.barcode.value = inventory.barcode || '';
      selectedAsset.sku.value = inventory.sku || '';
      selectedAsset.bought_at.value = inventory.bought_at || '';
      selectedAsset.return_location.value = inventory.return_location || '';
      selectedAsset.max_weight.value = inventory.max_weight || '';
      selectedAsset.min_weight.value = inventory.min_weight || '';
      selectedAsset.max_height.value = inventory.max_height || '';
      selectedAsset.min_height.value = inventory.min_height || '';
      selectedAsset.price.value = inventory.price || '';
      selectedAsset.quantity.value = inventory.quantity || '';
      selectedAsset.is_bookmarked.value = inventory.is_bookmarked || false;
      selectedAsset.is_returnable.value = inventory.is_returnable || Boolean(inventory.return_location) || false;
      selectedAsset.created_by.value = inventory.created_by || '';
      selectedAsset.created_at.value = inventory.created_at || '';
      selectedAsset.updated_by.value = inventory.updated_by || '';
      selectedAsset.updated_at.value = inventory.updated_at || '';
      selectedAsset.sharable_groups.value = inventory.sharable_groups || [];
      selectedAsset.creator_name = inventory?.creator_name || '';
      selectedAsset.updator_name = inventory?.updater_name || '';

      if (inventory?.return_datetime) {
        setReturnDateTime(dayjs(inventory.return_datetime));
      }

      if (inventory?.return_notes) {
        setOpenReturnNote(true);
        selectedAsset.return_notes.value = inventory.return_notes;
      }

      if (inventory?.color) {
        setColor(inventory.color);
      }

      setStorageLocation({ location: inventory.location });
      setFormData(selectedAsset);
    }
  }, [loading, inventory]);

  return (
    <>
      <RowHeader
        title="Viewing asset details"
        caption={`Editing ${formData.name.value}`}
        primaryStartIcon={<AddPhotoAlternateRounded />}
        primaryButtonTextLabel={'Add Image'}
        handleClickPrimaryButton={() => setEditImgMode(!editImgMode)}
      />
      <SelectedAssetFormFields
        formFields={formData}
        assetColor={color}
        handleColorChange={handleColorChange}
        selectedImage={selectedImage}
        handleInputChange={handleInputChange}
        options={storageLocations}
        storageLocation={storageLocation}
        setStorageLocation={setStorageLocation}
      />
      <Divider sx={{ marginTop: '1rem', marginBottom: '1rem' }}>
        <Typography variant="caption">More information</Typography>
      </Divider>
      <SelectedAssetMoreInformation
        formFields={formData}
        returnDateTime={returnDateTime}
        setReturnDateTime={setReturnDateTime}
        openReturnNote={openReturnNote}
        setOpenReturnNote={setOpenReturnNote}
        handleCheckbox={handleCheckbox}
        handleInputChange={handleInputChange}
      />

      <Divider sx={{ marginTop: '1rem', marginBottom: '1rem' }}>
        <Typography variant="caption">Weight and Dimension</Typography>
      </Divider>
      <SelectedAssetWeightDimension formFields={formData} handleInputChange={handleInputChange} />
      <Stack sx={{ margin: '1rem 0rem' }}>
        <Button
          startIcon={<CheckRounded fontSize="small" />}
          onClick={handleSubmit}
          disabled={formFieldChange ? isFormDisabled() : true}
          variant="outlined"
        >
          Submit
        </Button>
      </Stack>
      {editImgMode && (
        <SimpleModal
          title="Assign image"
          subtitle="Assign image to the selected item."
          handleClose={() => setEditImgMode(false)}
          maxSize="sm"
        >
          <ImagePicker id={id} name={formData.name.value} handleUpload={handleUpload} disableCancel />
        </SimpleModal>
      )}
    </>
  );
}
