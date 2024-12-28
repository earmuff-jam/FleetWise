import { AddRounded } from '@mui/icons-material';
import { IconButton, Stack, Typography } from '@mui/material';

export default function NoteHeader({ handleClick }) {
  return (
    <Stack direction="row" justifyContent="space-between">
      <Typography variant="h5" color="text.secondary" gutterBottom>
        Notes
      </Typography>
      <IconButton onClick={handleClick} size="small" data-tour="notes-1">
        <AddRounded fontSize="small" />
      </IconButton>
    </Stack>
  );
}
