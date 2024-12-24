import { profileActions } from '@features/Profile/profileSlice';
import Popper from '@material-ui/core/Popper/Popper';
import { ArrowDropDownRounded } from '@mui/icons-material';
import { Button, Grow, ButtonGroup, ClickAwayListener, MenuItem, MenuList, Paper } from '@mui/material';
import { useTour } from '@reactour/tour';
import { useRef, useState } from 'react';
import { useDispatch } from 'react-redux';

export default function AppToolbarButtonGroup({ profileDetails }) {
  const dispatch = useDispatch();
  const { setIsOpen } = useTour();

  const anchorRef = useRef(null);
  const [open, setOpen] = useState(false);

  const handleToggle = () => {
    setOpen((prevOpen) => !prevOpen);
  };

  const handleClose = (event) => {
    if (anchorRef.current && anchorRef.current.contains(event.target)) {
      return;
    }
    setOpen(false);
  };

  const handleAppearance = () => {
    const draftData = { ...profileDetails, appearance: !profileDetails.appearance || false };
    dispatch(profileActions.updateProfileDetails({ draftData }));
  };

  const options = [
    {
      id: 1,
      label: 'Help',
      action: () => setIsOpen(true),
    },
    {
      id: 2,
      label: 'Change Appearance',
      action: () => handleAppearance(),
    },
  ];

  return (
    <>
      <ButtonGroup variant="contained" ref={anchorRef} aria-label="button group to log off and help menu">
        <Button onClick={options[0].action}>{options[0].label}</Button>
        <Button
          size="small"
          aria-controls={open ? 'split-button-menu' : undefined}
          aria-expanded={open ? 'true' : undefined}
          aria-label="select more action"
          aria-haspopup="menu"
          onClick={handleToggle}
        >
          <ArrowDropDownRounded />
        </Button>
      </ButtonGroup>
      <Popper sx={{ zIndex: 1 }} open={open} anchorEl={anchorRef.current} role={undefined} transition disablePortal>
        {({ TransitionProps, placement }) => (
          <Grow
            {...TransitionProps}
            style={{
              transformOrigin: placement === 'bottom' ? 'center top' : 'center bottom',
            }}
          >
            <Paper>
              <ClickAwayListener onClickAway={handleClose}>
                <MenuList id="split-button-menu" autoFocusItem>
                  {options.map((option, index) => (
                    <MenuItem key={option.id} onClick={option.action}>
                      {option.label}
                    </MenuItem>
                  ))}
                </MenuList>
              </ClickAwayListener>
            </Paper>
          </Grow>
        )}
      </Popper>
    </>
  );
}
