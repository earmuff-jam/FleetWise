import { Paper, Skeleton, Table, TableContainer } from '@mui/material';

import { EmptyComponent } from '@common/utils';
import CustomTableBody from '@common/DataTable/CustomTableComponent/CustomTableBody/CustomTableBody';
import CustomTableHeader from '@common/DataTable/CustomTableComponent/CustomTableHeader/CustomTableHeader';

/**
 * TableComponent React Function - Displays the inventory table
 * @param {boolean} showActions - determines if actionButtons should be present, defaults: true
 * @param {boolean} hideCheckBox - determines if associated icons should be present, defaults: false
 * @param {boolean} hideIconButton - determines if associated icons should be present, defaults: false
 * @param {boolean} hideMoreDetailsButton - determines if associated icons should be present, defaults: false
 * @param {boolean} isLoading - determines if the selected data is still in loading state
 * @param {Array<Object>} columns - the columns to display for the table
 * @param {Function} rowFormatter - the row formatter to format each row
 * @param {Array<Object>} data - the data to display for each row in the table
 * @param {Array<String>} selectedIDList - the array of IDs that represent each item
 * @param {Function} onRowSelect - the function that is used to select a specific row
 * @param {Function} handleRowSelection - the function that is used to handle selection of rows
 * @param {Function} handleEdit - the function that is used to handle editing capabilities
 * @param {boolean} emptyComponentSubtext - subtitle text to display when there is no selected rows, defaults: empty string
 * @param {string} maxHeight - the maximum height of the table to expand, defaults: 20rem
 */
const TableComponent = ({
  paper = false,
  showActions = true,
  hideCheckBox = false,
  hideIconButton = false,
  hideMoreDetailsButton = false,
  isLoading,
  columns,
  data,
  rowFormatter,
  selectedIDList,
  onRowSelect,
  handleRowSelection,
  handleEdit,
  emptyComponentSubtext = '',
  maxHeight = '20rem',
}) => {
  if (isLoading) return <Skeleton height="10vh" />;

  if (!data || data.length === 0) {
    return <EmptyComponent subtitle={emptyComponentSubtext} />;
  }

  return (
    <TableContainer component={paper ? Paper : null} sx={{ maxHeight: maxHeight }}>
      <Table>
        <CustomTableHeader
          columns={columns}
          selectedIDList={selectedIDList}
          showActions={showActions}
          hideCheckBox={hideCheckBox}
          handleRowSelection={handleRowSelection}
        />
        <CustomTableBody
          data={data}
          columns={columns}
          selectedIDList={selectedIDList}
          hideCheckBox={hideCheckBox}
          handleRowSelection={handleRowSelection}
          rowFormatter={rowFormatter}
          showActions={showActions}
          handleEdit={handleEdit}
          onRowSelect={onRowSelect}
          hideIconButton={hideIconButton}
          hideMoreDetailsButton={hideMoreDetailsButton}
        />
      </Table>
    </TableContainer>
  );
};

export default TableComponent;
