import { useEffect, useState } from 'react';

import dayjs from 'dayjs';
import { enqueueSnackbar } from 'notistack';
import { useDispatch, useSelector } from 'react-redux';
import { Box, Skeleton, Typography } from '@mui/material';
import { inventoryActions } from '@features/Assets/inventorySlice';

import AddAssetSteps from '@features/Assets/AddAssetFormDetails/AddAssetSteps';
import { ADD_ASSET_FORM } from '@features/Assets/AddAssetFormDetails/constants';
import AddAssetFormSelection from '@features/Assets/AddAssetFormDetails/AddAssetFormSelection';
import AddAssetActionButtons from '@features/Assets/AddAssetFormDetails/AddAssetActionButtons';
import AddAssetFormInstructions from '@features/Assets/AddAssetFormDetails/AddAssetFormInstructions';

export default function AddAssetDetails({ handleClose }) {
  const dispatch = useDispatch();
  const { storageLocations: options, isLoading } = useSelector((state) => state.inventory);

  const [activeStep, setActiveStep] = useState(0);
  const [skipped, setSkipped] = useState(new Set());
  const [storageLocation, setStorageLocation] = useState('');
  const [returnDateTime, setReturnDateTime] = useState(null);

  const [formData, setFormData] = useState({ ...ADD_ASSET_FORM });

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

    setFormData(updatedFormData);
  };

  const handleCheckbox = (name, value) => {
    setFormData((prevFormData) => ({
      ...prevFormData,
      [name]: { ...prevFormData[name], value },
    }));
  };

  const handleSubmit = (event) => {
    event.preventDefault();
    const containsErr = Object.values(formData).reduce((acc, el) => {
      if (el.errorMsg) {
        return true;
      }
      return acc;
    }, false);

    const requiredFormFields = Object.values(formData).filter((v) => v.required);
    const isRequiredFieldsEmpty = requiredFormFields.some((el) => el.value.trim() === '');

    if (containsErr || isRequiredFieldsEmpty || storageLocation === null || Object.keys(storageLocation).length <= 0) {
      enqueueSnackbar('Unable to save new asset. Missing required details.', {
        variant: 'error',
      });
      return;
    }

    const formattedData = Object.values(formData).reduce((acc, el) => {
      if (el.id === 'quantity') {
        acc[el.id] = parseInt(el.value);
      } else if (el.id === 'price') {
        acc[el.id] = parseFloat(el.value);
      } else if (el.value) {
        acc[el.id] = el.value;
      }
      return acc;
    }, {});

    const storageLocationID =
      options.find((v) => v.location === storageLocation.location)?.id || storageLocation.location;

    const draftRequest = {
      ...formattedData,
      location: storageLocationID,
      return_datetime: returnDateTime !== null ? returnDateTime.toISOString() : null,
      created_at: dayjs().toISOString(),
    };
    dispatch(inventoryActions.addInventory(draftRequest));
    setFormData({ ...ADD_ASSET_FORM });
    enqueueSnackbar('Added new asset.', {
      variant: 'success',
    });
    handleClose();
  };

  const isStepOptional = (step) => {
    return step === 1;
  };

  const isStepSkipped = (step) => {
    return skipped.has(step);
  };

  const handleNext = () => {
    let newSkipped = skipped;
    if (isStepSkipped(activeStep)) {
      newSkipped = new Set(newSkipped.values());
      newSkipped.delete(activeStep);
    }

    setActiveStep((prevActiveStep) => prevActiveStep + 1);
    setSkipped(newSkipped);
  };

  const handleBack = () => {
    setActiveStep((prevActiveStep) => prevActiveStep - 1);
  };

  const handleSkip = () => {
    if (!isStepOptional(activeStep)) {
      throw new Error("You can't skip a step that isn't optional.");
    }

    setActiveStep((prevActiveStep) => prevActiveStep + 1);
    setSkipped((prevSkipped) => {
      const newSkipped = new Set(prevSkipped.values());
      newSkipped.add(activeStep);
      return newSkipped;
    });
  };

  const handleReset = () => {
    setActiveStep(0);
  };

  const isDisabled = () => {
    const containsErr = Object.values(formData).reduce((acc, el) => {
      if (el.errorMsg) {
        return true;
      }
      return acc;
    }, false);

    const requiredFormFields = Object.values(formData).filter((v) => v.isRequired);
    const isRequiredFieldsEmpty = requiredFormFields.some((el) => el.value.trim() === '');

    return containsErr || isRequiredFieldsEmpty || storageLocation === null || Object.keys(storageLocation).length <= 0;
  };

  useEffect(() => {
    dispatch(inventoryActions.getStorageLocations());
  }, []);

  if (isLoading) return <Skeleton height="30vh" />;

  return (
    <Box sx={{ width: '100%' }}>
      <AddAssetSteps activeStep={activeStep} isStepOptional={isStepOptional} isStepSkipped={isStepSkipped} />
      <Typography sx={{ mt: 2, mb: 1 }}>Step {activeStep + 1}</Typography>
      <AddAssetFormInstructions stepNumber={activeStep + 1} />
      <AddAssetFormSelection
        stepNumber={activeStep + 1}
        formData={formData}
        storageLocation={storageLocation}
        setStorageLocation={setStorageLocation}
        handleInputChange={handleInputChange}
        handleCheckbox={handleCheckbox}
        handleReset={handleReset}
        handleSubmit={handleSubmit}
        options={options}
        returnDateTime={returnDateTime}
        setReturnDateTime={setReturnDateTime}
      />
      <AddAssetActionButtons
        activeStep={activeStep}
        handleBack={handleBack}
        handleNext={handleNext}
        handleSkip={handleSkip}
        isStepOptional={isStepOptional}
        isNextDisabled={isDisabled()}
        isBackDisabled={activeStep === 0}
      />
    </Box>
  );
}
