import { IconButton, Stack, TableCell } from '@mui/material';
import { EditRounded, FileOpenRounded } from '@mui/icons-material';

export default function CustomShowActionsTableCell({
  row,
  selectedID,
  hideIconButton,
  hideMoreDetailsButton,
  handleEdit,
  onRowSelect,
}) {
  return (
    <TableCell>
      <Stack direction="row">
        {!hideIconButton && (
          <IconButton onClick={() => handleEdit(selectedID)} size="small" data-tour="12">
            <EditRounded color="primary" fontSize="small" />
          </IconButton>
        )}
        {!hideMoreDetailsButton && (
          <IconButton size="small" onClick={() => onRowSelect(row)} data-tour="11">
            <FileOpenRounded color="primary" fontSize="small" />
          </IconButton>
        )}
      </Stack>
    </TableCell>
  );
}
