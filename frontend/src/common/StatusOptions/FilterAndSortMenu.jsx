import { useState } from 'react';

import { Box, IconButton, Menu, MenuItem, Paper, Stack, Tooltip, Typography } from '@mui/material';

import { CloseRounded, FilterAltRounded, SortRounded } from '@mui/icons-material';

import { STATUS_OPTIONS } from '@common/StatusOptions/constants';

export default function FilterAndSortMenu({
  sortingOrder,
  setSortingOrder,
  selectedFilter,
  setSelectedFilter,
  filterBtnDataTour,
  sortBtnDataTour,
}) {
  const [anchorEl, setAnchorEl] = useState(null);
  const open = Boolean(anchorEl);

  const handleOpenMenu = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuItemClick = (optionLabel) => {
    setSelectedFilter(optionLabel);
    setAnchorEl(null);
  };

  const handleCloseMenu = () => {
    setAnchorEl(null);
  };

  const toggleSort = () => setSortingOrder(!sortingOrder);

  return (
    <Stack sx={{ flexDirection: 'row', gap: '1rem', mb: '1rem' }}>
      <IconButton size="small" onClick={handleOpenMenu} data-tour={filterBtnDataTour}>
        <FilterAltRounded fontSize="small" />
      </IconButton>
      <Tooltip title={`${sortingOrder ? 'DESC' : 'ASC'}`}>
        <IconButton size="small" onClick={toggleSort} disabled={Boolean(selectedFilter)} data-tour={sortBtnDataTour}>
          <SortRounded fontSize="small" />
        </IconButton>
      </Tooltip>
      {Boolean(selectedFilter) && (
        <IconButton size="small" onClick={() => setSelectedFilter('')}>
          <CloseRounded fontSize="small" color="error" />
        </IconButton>
      )}
      <Menu
        id="filter-sort-menu"
        anchorEl={anchorEl}
        open={open}
        onClose={handleCloseMenu}
        slotProps={{
          root: { sx: { '.MuiList-root': { padding: 0 } } },
        }}
      >
        <Paper sx={{ width: 200 }}>
          <Typography variant="h6" sx={{ paddingLeft: '0.5rem', paddingTop: '0.5rem' }}>
            Filter results:
          </Typography>
          {STATUS_OPTIONS.map((option, index) => (
            <MenuItem
              key={index}
              id={option.label}
              selected={option.label === selectedFilter}
              onClick={() => handleMenuItemClick(option.label)}
            >
              <Stack direction="row" alignItems="center" justifyContent="space-between" width="100%">
                <Box>{option.display}</Box>
                <Box component="span" sx={{ display: 'flex', alignItems: 'center' }}>
                  {option.icon}
                </Box>
              </Stack>
            </MenuItem>
          ))}
        </Paper>
      </Menu>
    </Stack>
  );
}
