import { Paper, Stack } from '@mui/material';
import RowHeader from '@common/RowHeader';
import { pluralizeWord } from '@common/utils';
import { AddRounded, RemoveRounded } from '@mui/icons-material';
import { VIEW_INVENTORY_LIST_HEADERS } from '@features/InventoryList/constants';
import TableComponent from '@common/DataTable/CustomTableComponent/TableComponent';

export default function MaintenancePlanItemDetailsContent({
  selectedIDList,
  setSelectedIDList,
  itemsInMaintenancePlan,
  handleOpenModal,
  handleRemoveAssociation,
}) {
  const handleRowSelection = (_, id) => {
    if (id === 'all') {
      if (selectedIDList.length !== 0) {
        setSelectedIDList([]);
      } else {
        setSelectedIDList(itemsInMaintenancePlan.map((v) => v.id));
      }
    } else {
      const selectedIndex = selectedIDList.indexOf(id);
      let draftSelected = [];
      if (selectedIndex === -1) {
        draftSelected = draftSelected.concat(selectedIDList, id);
      } else if (selectedIndex === 0) {
        draftSelected = draftSelected.concat(selectedIDList.slice(1));
      } else if (selectedIndex === selectedIDList.length - 1) {
        draftSelected = draftSelected.concat(selectedIDList.slice(0, -1));
      } else if (selectedIndex > 0) {
        draftSelected = draftSelected.concat(
          selectedIDList.slice(0, selectedIndex),
          selectedIDList.slice(selectedIndex + 1)
        );
      }
      setSelectedIDList(draftSelected);
    }
  };

  const rowFormatter = (row, columnName, columnData) => {
    if (columnData.modifier) {
      return columnData.modifier(row[columnName] || '-');
    } else {
      return row[columnName] || '-';
    }
  };

  return (
    <Paper elevation={1} sx={{ padding: '1rem' }}>
      <Stack spacing={2}>
        <RowHeader
          title="Items"
          caption={`Total ${pluralizeWord('item', itemsInMaintenancePlan?.length || 0)}`}
          primaryButtonTextLabel="Add"
          primaryStartIcon={<AddRounded />}
          handleClickPrimaryButton={handleOpenModal}
          secondaryButtonTextLabel="Remove"
          secondaryStartIcon={<RemoveRounded color="error" />}
          handleClickSecondaryButton={handleRemoveAssociation}
          secondaryButtonDisabled={selectedIDList.length <= 0}
        />
        <TableComponent
          showActions={false}
          data={itemsInMaintenancePlan}
          columns={Object.values(VIEW_INVENTORY_LIST_HEADERS).filter((v) => v.displayConcise)}
          rowFormatter={rowFormatter}
          selectedIDList={selectedIDList}
          handleRowSelection={handleRowSelection}
          emptyComponentSubtext="Associate assets."
        />
      </Stack>
    </Paper>
  );
}
