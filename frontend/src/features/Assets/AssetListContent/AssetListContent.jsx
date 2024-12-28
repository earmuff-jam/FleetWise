import { useState } from 'react';

import { useDispatch } from 'react-redux';
import { useNavigate } from 'react-router-dom';

import SimpleModal from '@common/SimpleModal';
import { Stack, TextField } from '@mui/material';
import GridComponent from '@features/Assets/AssetListContent/GridComponent';
import { inventoryActions } from '@features/Assets/inventorySlice';

import { CheckRounded, CloseRounded, EditRounded } from '@mui/icons-material';
import AddAssetsInBulk from '@features/Assets/AddAssetsInBulk/AddAssetsInBulk';
import AddAssetDetails from '@features/Assets/AddAssetFormDetails/AddAssetDetails';

import TableComponent from '@common/DataTable/CustomTableComponent/TableComponent';
import { MODAL_STATE, VIEW_INVENTORY_LIST_HEADERS } from '@features/Assets/constants';
import AssetDetailsDrawer from '@features/Assets/AssetDetailsDrawer/AssetDetailsDrawer';

export default function AssetListContent({
  loading,
  modalState,
  setModalState,
  inventories,
  options,
  gridMode,
  rowSelected,
  setRowSelected,
  handleCloseModal,
}) {
  const navigate = useNavigate();
  const dispatch = useDispatch();

  const [inputColumn, setInputColumn] = useState('');
  const [selectedRow, setSelectedRow] = useState([]); // to display more details
  const [editLineItem, setEditLineItem] = useState({ editItem: false, rowID: -1, column: '' });

  const onRowSelect = (row) => {
    setModalState(MODAL_STATE.MORE_DETAILS);
    setSelectedRow(row);
  };

  const handleEdit = (itemID) => {
    navigate(`/inventories/${itemID}/update`);
  };

  const updateSelectedCol = (rowID, columnName, inputColumn) => {
    dispatch(inventoryActions.updateAssetCol({ assetID: rowID, columnName, inputColumn }));
  };

  // checkbox actions
  const handleRowSelection = (_, id) => {
    if (id === 'all') {
      if (rowSelected.length === 0) {
        setRowSelected(inventories.map((v) => v.id));
      } else {
        setRowSelected([]);
      }
    } else {
      const selectedIndex = rowSelected.indexOf(id);
      let draftSelected = [];
      if (selectedIndex === -1) {
        draftSelected = draftSelected.concat(rowSelected, id);
      } else if (selectedIndex === 0) {
        draftSelected = draftSelected.concat(rowSelected.slice(1));
      } else if (selectedIndex === rowSelected.length - 1) {
        draftSelected = draftSelected.concat(rowSelected.slice(0, -1));
      } else if (selectedIndex > 0) {
        draftSelected = draftSelected.concat(rowSelected.slice(0, selectedIndex), rowSelected.slice(selectedIndex + 1));
      }
      setRowSelected(draftSelected);
    }
  };

  const populateIcon = (editLineItem, row, column, inputColumn, setInputColumn) => {
    if (editLineItem.editItem && editLineItem.rowID === row.id && editLineItem.column === column && inputColumn) {
      return (
        <CheckRounded
          sx={{ height: '1rem', width: '1rem', marginLeft: '0.5rem', cursor: 'pointer' }}
          color="primary"
          onClick={() => {
            updateSelectedCol(row.id, editLineItem.column, inputColumn);
            setInputColumn('');
            setEditLineItem({ editItem: false, rowID: -1, column: '' });
          }}
        />
      );
    } else if (editLineItem.editItem && editLineItem.rowID === row.id && editLineItem.column === column) {
      return (
        <CloseRounded
          sx={{ height: '1rem', width: '1rem', marginLeft: '0.5rem', cursor: 'pointer' }}
          onClick={() => setEditLineItem({ editItem: false, rowID: -1, column: '' })}
        />
      );
    } else {
      return (
        <EditRounded
          sx={{ height: '1rem', width: '1rem', marginLeft: '0.5rem', cursor: 'pointer' }}
          onClick={() => setEditLineItem({ editItem: true, rowID: row.id, column: column })}
        />
      );
    }
  };

  const rowFormatter = (row, columnName, columnData) => {
    const handleInputChange = (ev) => {
      const value = ev.target.value;
      if (!(isNaN(value) || parseInt(value) <= 0)) {
        // only +ve numbers
        setInputColumn(value);
      }
    };

    if (['price', 'quantity'].includes(columnName)) {
      return (
        <Stack direction="row" alignItems="center">
          {editLineItem.editItem && editLineItem.rowID === row.id && editLineItem.column === columnName ? (
            <TextField variant="standard" label={row[columnName]} value={inputColumn} onChange={handleInputChange} />
          ) : row[columnName] <= 0 ? (
            '-'
          ) : (
            row[columnName]
          )}
          {populateIcon(editLineItem, row, columnName, inputColumn, setInputColumn)}
        </Stack>
      );
    } else if (columnData.modifier) {
      return columnData.modifier(row[columnName] || '-');
    } else {
      return row[columnName] || '-';
    }
  };

  return (
    <>
      <Stack>
        {gridMode ? (
          <GridComponent
            isLoading={loading}
            data={options}
            rowSelected={rowSelected}
            handleRowSelection={handleRowSelection}
          />
        ) : (
          <TableComponent
            paper
            isLoading={loading}
            data={options}
            columns={Object.values(VIEW_INVENTORY_LIST_HEADERS).filter((v) => v.displayConcise)}
            rowFormatter={rowFormatter}
            selectedIDList={rowSelected}
            onRowSelect={onRowSelect}
            handleRowSelection={handleRowSelection}
            handleEdit={handleEdit}
            maxHeight={'100%'}
          />
        )}
      </Stack>
      {modalState === MODAL_STATE.ADD_ITEM && (
        <SimpleModal title="Add New Item" handleClose={handleCloseModal} maxSize="sm">
          <AddAssetDetails handleClose={handleCloseModal} />
        </SimpleModal>
      )}
      {modalState === MODAL_STATE.BULK_ITEM && (
        <SimpleModal
          title="Add Bulk Item"
          subtitle="Upload inventory items in bulk with the selected template."
          handleClose={handleCloseModal}
          maxSize="sm"
        >
          <AddAssetsInBulk handleClose={handleCloseModal} />
        </SimpleModal>
      )}
      {modalState === MODAL_STATE.MORE_DETAILS && (
        <AssetDetailsDrawer
          columns={Object.values(VIEW_INVENTORY_LIST_HEADERS)}
          resetSelection={handleCloseModal}
          title="View Item Details"
          selectedRow={selectedRow}
        />
      )}
    </>
  );
}
