import { useState } from 'react';

import { Autocomplete, IconButton, Stack, TextField } from '@mui/material';

import { GridViewRounded, SearchRounded, ViewListRounded } from '@mui/icons-material';

import RowHeader from '@common/RowHeader';
import { pluralizeWord } from '@common/utils';
import { MODAL_STATE } from '@features/Assets/constants';
import VerticalMenu from '@features/Assets/AssetListHeader/VerticalMenu';

export default function AssetListHeader({
  gridMode,
  setGridMode,
  setModalState,
  inventories,
  options,
  setOptions,
  disableDelete,
  handleRemoveInventory,
}) {
  const [showSearch, setShowSearch] = useState(false);

  return (
    <Stack direction="row" justifyContent="space-between">
      <RowHeader title="Assets" caption={`Viewing ${pluralizeWord('asset detail', options?.length)}`} />
      <Stack direction="row" alignItems="center">
        {showSearch ? (
          <Autocomplete
            sx={{ minWidth: '15rem' }}
            id="inventory-items-autocomplete"
            options={options}
            autoHighlight
            forcePopupIcon
            getOptionLabel={(option) => option.name}
            onChange={(_, newValue) => {
              if (newValue) {
                setOptions(inventories.filter((option) => option.id === newValue.id));
              } else {
                setOptions(inventories);
              }
            }}
            renderInput={(params) => <TextField variant="standard" {...params} label="Search ..." />}
          />
        ) : (
          <IconButton onClick={() => setShowSearch(!showSearch)} data-tour="assets-1">
            <SearchRounded />
          </IconButton>
        )}
        <IconButton size="small" onClick={() => setGridMode(!gridMode)} data-tour="assets-2">
          {!gridMode ? (
            <GridViewRounded color="primary" fontSize="small" />
          ) : (
            <ViewListRounded color="primary" fontSize="small" />
          )}
        </IconButton>
        <VerticalMenu
          disableDelete={disableDelete}
          handleAddInventory={() => setModalState(MODAL_STATE.ADD_ITEM)}
          handleBulkInventory={() => setModalState(MODAL_STATE.BULK_ITEM)}
          handleRemoveInventory={handleRemoveInventory}
        />
      </Stack>
    </Stack>
  );
}
