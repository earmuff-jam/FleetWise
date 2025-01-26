import { useState } from 'react';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import { IconButton, Menu, MenuItem, Typography } from '@mui/material';

export default function VerticalMenu({
  disableDelete,
  handleAddInventory,
  handleBulkInventory,
  handleRemoveInventory,
}) {
  const [anchorEl, setAnchorEl] = useState(null);
  const open = Boolean(anchorEl);

  const handleClick = (event) => setAnchorEl(event.currentTarget);
  const handleClose = () => setAnchorEl(null);

  const options = [
    { id: 1, label: 'Add New', action: handleAddInventory, disabled: false },
    { id: 2, label: 'Bulk Add', action: handleBulkInventory, disabled: false },
    { id: 3, label: 'Delete Chosen', action: handleRemoveInventory, disabled: disableDelete },
  ];

  return (
    <>
      <IconButton
        aria-label="more"
        id="long-button"
        aria-controls={open ? 'long-menu' : undefined}
        aria-expanded={open ? 'true' : undefined}
        aria-haspopup="true"
        onClick={handleClick}
        data-tour="assets-3"
      >
        <MoreVertIcon />
      </IconButton>
      <Menu
        id="long-menu"
        MenuListProps={{
          'aria-labelledby': 'long-button',
        }}
        anchorEl={anchorEl}
        open={open}
        onClose={handleClose}
      >
        {options.map((option) => (
          <MenuItem
            key={option.id}
            disabled={option.disabled}
            onClick={() => {
              option.action();
              handleClose();
            }}
          >
            <Typography variant="caption">{option.label}</Typography>
          </MenuItem>
        ))}
      </Menu>
    </>
  );
}
