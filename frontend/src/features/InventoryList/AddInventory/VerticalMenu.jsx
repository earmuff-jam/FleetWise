import { IconButton, Menu, MenuItem } from '@mui/material';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import { useState } from 'react';

export default function VerticalMenu({ rowSelected, handleAddInventory, handleBulkInventory, handleRemoveInventory }) {
  const [anchorEl, setAnchorEl] = useState(null);
  const open = Boolean(anchorEl);

  const handleClick = (event) => setAnchorEl(event.currentTarget);
  const handleClose = () => setAnchorEl(null);

  const options = [
    { id: 1, label: 'Add Item', action: handleAddInventory, disabled: false },
    { id: 2, label: 'Add Bulk', action: handleBulkInventory, disabled: false },
    { id: 3, label: 'Delete Selected', action: handleRemoveInventory, disabled: rowSelected.length <= 0 },
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
            {option.label}
          </MenuItem>
        ))}
      </Menu>
    </>
  );
}
