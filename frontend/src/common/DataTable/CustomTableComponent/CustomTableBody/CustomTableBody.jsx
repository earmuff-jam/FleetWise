import { TableBody, TableCell, TableRow } from '@mui/material';
import CustomCheckboxTableCell from '@common/DataTable/CustomTableComponent/CustomTableBody/CustomCheckboxTableCell';
import CustomShowActionsTableCell from '@common/DataTable/CustomTableComponent/CustomTableBody/CustomShowActionsTableCell';

export default function CustomTableBody({
  data,
  columns,
  selectedIDList,
  hideCheckBox,
  handleRowSelection,
  rowFormatter,
  showActions,
  handleEdit,
  onRowSelect,
  hideIconButton,
  hideMoreDetailsButton,
}) {
  return (
    <TableBody>
      {data.map((row) => {
        const isSelected = (id) => selectedIDList && selectedIDList.indexOf(id) !== -1;
        const selectedID = row.id;
        const isItemSelected = isSelected(selectedID);
        return (
          <TableRow hover key={row.id}>
            {!hideCheckBox && (
              <CustomCheckboxTableCell
                selectedID={selectedID}
                isItemSelected={isItemSelected}
                handleRowSelection={handleRowSelection}
              />
            )}
            {Object.keys(columns).map((colKey) => {
              const column = columns[colKey];
              return (
                <TableCell key={column.id} size="small">
                  {rowFormatter(row, column.colName, column)}
                </TableCell>
              );
            })}
            {showActions && (
              <CustomShowActionsTableCell
                row={row}
                selectedID={selectedID}
                handleEdit={handleEdit}
                onRowSelect={onRowSelect}
                hideIconButton={hideIconButton}
                hideMoreDetailsButton={hideMoreDetailsButton}
              />
            )}
          </TableRow>
        );
      })}
    </TableBody>
  );
}
