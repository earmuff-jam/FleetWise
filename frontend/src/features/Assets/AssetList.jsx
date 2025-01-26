import { useEffect, useState } from 'react';

import { useDispatch, useSelector } from 'react-redux';

import { Stack } from '@mui/material';
import { ConfirmationBoxModal } from '@common/utils';
import { MODAL_STATE } from '@features/Assets/constants';
import { inventoryActions } from '@features/Assets/inventorySlice';
import AssetListHeader from '@features/Assets/AssetListHeader/AssetListHeader';
import AssetListContent from '@features/Assets/AssetListContent/AssetListContent';

export default function AssetList() {
  const dispatch = useDispatch();
  const { loading, inventories = [] } = useSelector((state) => state.inventory);

  const [options, setOptions] = useState([]);
  const [rowSelected, setRowSelected] = useState([]); // this is for checkbox and associated actions
  const [gridMode, setGridMode] = useState(false);

  const [openDialog, setOpenDialog] = useState(false);
  const [idToDelete, setIdToDelete] = useState(-1);
  const [modalState, setModalState] = useState(MODAL_STATE.NONE);

  const handleCloseModal = () => setModalState(MODAL_STATE.NONE);

  const handleRemoveInventory = () => {
    setOpenDialog(true);
    setIdToDelete(rowSelected);
  };

  const reset = () => {
    setOpenDialog(false);
    setRowSelected([]);
    setIdToDelete(-1);
  };

  const confirmDelete = (id) => {
    if (id === -1) {
      return;
    }
    dispatch(inventoryActions.removeInventoryRows(rowSelected));
    reset();
  };

  useEffect(() => {
    if (Array.isArray(inventories)) {
      setOptions(inventories);
    }
  }, [loading]);

  useEffect(() => {
    dispatch(inventoryActions.getAllInventoriesForUser());
  }, []);

  return (
    <Stack flexGrow="1" spacing={2} data-tour="assets-0">
      <AssetListHeader
        gridMode={gridMode}
        setGridMode={setGridMode}
        options={options}
        setOptions={setOptions}
        inventories={inventories}
        setModalState={setModalState}
        handleRemoveInventory={handleRemoveInventory}
        disableDelete={rowSelected.length <= 0}
      />
      <AssetListContent
        loading={loading}
        modalState={modalState}
        setModalState={setModalState}
        gridMode={gridMode}
        inventories={inventories}
        options={options}
        rowSelected={rowSelected}
        setRowSelected={setRowSelected}
        handleCloseModal={handleCloseModal}
      />
      <ConfirmationBoxModal
        openDialog={openDialog}
        title="Confirm deletion"
        handleClose={reset}
        maxSize="xs"
        deleteID={idToDelete}
        confirmDelete={confirmDelete}
      />
    </Stack>
  );
}
