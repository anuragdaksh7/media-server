// @ts-nocheck;

import Icon from './Icon';

type Props = {
  name?: string;
  size?: number;
  color?: string;
  stroke?: string;
  onClick?: () => void
  cursor?: string;
  className?: string;
};

export const Logo = (props: Props) => <Icon {...props} name="logo" />;