import { Button, IconButton, Stack } from '@mui/material';
import { AddRounded, FileDownload } from '@mui/icons-material';

export default function SectionCardHeaderButton({
  title,
  handleButtonClick,
  handleIconButtonClick,
  disableDownloadIcon,
}) {
  return (
    <Stack direction="row" spacing={1}>
      <Button onClick={handleButtonClick} startIcon={<AddRounded />} variant="outlined">
        {title}
      </Button>
      <IconButton size="small" onClick={handleIconButtonClick} disabled={disableDownloadIcon}>
        <FileDownload fontSize="small" />
      </IconButton>
    </Stack>
  );
}
