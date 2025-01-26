import { useDispatch } from 'react-redux';
import { Autocomplete, CardMedia, Stack, TextField, createFilterOptions } from '@mui/material';

import { inventoryActions } from '@features/Assets/inventorySlice';
import TextFieldWithLabel from '@common/TextFieldWithLabel/TextFieldWithLabel';
import ColorPicker from '@common/ColorPicker';

const filter = createFilterOptions();

export default function SelectedAssetFormFields({
  formFields,
  selectedImage,
  assetColor,
  handleColorChange,
  handleInputChange,
  options,
  storageLocation,
  setStorageLocation,
}) {
  const dispatch = useDispatch();
  return (
    <Stack spacing={2} marginTop={'1rem'}>
      <Stack direction="row" spacing={2} alignItems={'center'}>
        <CardMedia
          component="img"
          sx={{ width: 300 }}
          image={selectedImage || '/blank_canvas.png'}
          alt="Default placeholder image for the asset"
        />
        <Stack spacing={2} flexGrow={1}>
          <TextFieldWithLabel
            id={formFields.name.name}
            name={formFields.name.name}
            label={formFields.name.label}
            value={formFields.name.value}
            size={formFields.name.size}
            placeholder={formFields.name.placeholder}
            handleChange={handleInputChange}
            required={formFields.name.required}
            fullWidth={formFields.name.fullWidth}
            error={Boolean(formFields.name.errorMsg)}
            helperText={formFields.name.errorMsg}
            variant={formFields.name.variant}
          />
          <TextFieldWithLabel
            id={formFields.description.name}
            name={formFields.description.name}
            label={formFields.description.label}
            value={formFields.description.value}
            size={formFields.description.size}
            placeholder={formFields.description.placeholder}
            handleChange={handleInputChange}
            required={formFields.description.required}
            fullWidth={formFields.description.fullWidth}
            error={Boolean(formFields.description.errorMsg)}
            helperText={formFields.description.errorMsg}
            variant={formFields.description.variant}
            multiline={true}
            rows={4}
          />
        </Stack>
      </Stack>
      <Stack direction="row" spacing={2}>
        <TextFieldWithLabel
          id={formFields.price.name}
          name={formFields.price.name}
          label={formFields.price.label}
          value={formFields.price.value}
          size={formFields.price.size}
          placeholder={formFields.price.placeholder}
          handleChange={handleInputChange}
          required={formFields.price.required}
          fullWidth={formFields.price.fullWidth}
          error={Boolean(formFields.price.errorMsg)}
          helperText={formFields.price.errorMsg}
          variant={formFields.price.variant}
        />
        <TextFieldWithLabel
          id={formFields.barcode.name}
          name={formFields.barcode.name}
          label={formFields.barcode.label}
          value={formFields.barcode.value}
          size={formFields.barcode.size}
          placeholder={formFields.barcode.placeholder}
          handleChange={handleInputChange}
          required={formFields.barcode.required}
          fullWidth={formFields.barcode.fullWidth}
          error={Boolean(formFields.barcode.errorMsg)}
          helperText={formFields.barcode.errorMsg}
          variant={formFields.barcode.variant}
        />
        <TextFieldWithLabel
          id={formFields.sku.name}
          name={formFields.sku.name}
          label={formFields.sku.label}
          value={formFields.sku.value}
          size={formFields.sku.size}
          placeholder={formFields.sku.placeholder}
          handleChange={handleInputChange}
          required={formFields.sku.required}
          fullWidth={formFields.sku.fullWidth}
          error={Boolean(formFields.sku.errorMsg)}
          helperText={formFields.sku.errorMsg}
          variant={formFields.sku.variant}
        />
      </Stack>
      <Stack direction="row" spacing={1}>
        <TextFieldWithLabel
          id={formFields.quantity.name}
          name={formFields.quantity.name}
          label={formFields.quantity.label}
          value={formFields.quantity.value}
          size={formFields.quantity.size}
          placeholder={formFields.quantity.placeholder}
          handleChange={handleInputChange}
          required={formFields.quantity.required}
          fullWidth={formFields.quantity.fullWidth}
          error={Boolean(formFields.quantity.errorMsg)}
          helperText={formFields.quantity.errorMsg}
          variant={formFields.quantity.variant}
        />
        <ColorPicker value={assetColor} handleChange={handleColorChange} label={'Associate color'} />
        <TextFieldWithLabel
          id={formFields.bought_at.name}
          name={formFields.bought_at.name}
          label={formFields.bought_at.label}
          value={formFields.bought_at.value}
          size={formFields.bought_at.size}
          placeholder={formFields.bought_at.placeholder}
          handleChange={handleInputChange}
          required={formFields.bought_at.required}
          fullWidth={formFields.bought_at.fullWidth}
          error={Boolean(formFields.bought_at.errorMsg)}
          helperText={formFields.bought_at.errorMsg}
          variant={formFields.bought_at.variant}
        />
      </Stack>
      <Autocomplete
        fullWidth
        freeSolo
        forcePopupIcon
        value={storageLocation || ''}
        onOpen={() => dispatch(inventoryActions.getStorageLocations())}
        onChange={(event, newValue) => {
          if (typeof newValue === 'string') {
            setStorageLocation({
              location: newValue,
            });
          } else if (newValue && newValue.inputValue) {
            setStorageLocation({
              location: newValue.inputValue,
            });
          } else {
            setStorageLocation(newValue);
          }
        }}
        filterOptions={(options, params) => {
          const filtered = filter(options, params);
          const { inputValue } = params;
          // Suggest the creation of a new value
          const isExisting = options.some((option) => inputValue === option.location);
          if (inputValue !== '' && !isExisting) {
            filtered.push({
              inputValue,
              location: `Add "${inputValue}"`,
            });
          }
          return filtered;
        }}
        selectOnFocus
        clearOnBlur
        handleHomeEndKeys
        id="autocomplete-storage-location"
        options={options || []}
        getOptionLabel={(option) => {
          if (typeof option === 'string') {
            return option;
          }
          if (option.inputValue) {
            return option.inputValue;
          }
          return option.location || '';
        }}
        renderOption={(props, option) => (
          <li {...props} key={option.location}>
            {option.location}
          </li>
        )}
        renderInput={(params) => (
          <TextField
            {...params}
            fullWidth
            label="Storage Location"
            variant="standard"
            placeholder="Where do you plan to store this item"
          />
        )}
      />
    </Stack>
  );
}
