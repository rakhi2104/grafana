import React, { FC } from 'react';
import { css } from 'emotion';
import tinycolor from 'tinycolor2';
import { getTagColorsFromName } from '../../utils';
import { stylesFactory, useTheme } from '../../themes';
import { Icon } from '../Icon/Icon';
import { GrafanaTheme } from '@grafana/data';

interface Props {
  name: string;

  onRemove: (tag: string) => void;
}

const getStyles = stylesFactory(({ theme, name }: { theme: GrafanaTheme; name: string }) => {
  const { color } = getTagColorsFromName(name);
  let hoverColor;
  if (theme.isDark) {
    hoverColor = tinycolor(color)
      .darken()
      .toHexString();
  } else {
    hoverColor = tinycolor(color)
      .lighten()
      .toHexString();
  }
  return {
    itemStyle: css`
      font-weight: ${theme.typography.weight.bold};
      font-size: ${theme.typography.size.sm};
      line-height: 16px;
      vertical-align: baseline;
      background-color: ${color};
      color: ${theme.palette.gray98};
      white-space: nowrap;
      text-shadow: none;
      padding: 4px 8px;
      border-radius: ${theme.border.radius.sm};

      :hover {
        background-color: ${hoverColor};
        cursor: pointer;
      }
    `,

    nameStyle: css`
      margin-right: 3px;
    `,

    iconStyle: css`
      margin-bottom: 0;
    `,
  };
});

export const TagItem: FC<Props> = ({ name, onRemove }) => {
  const theme = useTheme();
  const styles = getStyles({ theme, name });

  return (
    <div className={styles.itemStyle}>
      <span className={styles.nameStyle}>{name}</span>
      <Icon className={styles.iconStyle} name="times" onClick={() => onRemove(name)} />
    </div>
  );
};
