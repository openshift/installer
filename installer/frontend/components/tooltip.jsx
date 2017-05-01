import React from 'react';

import _ from 'lodash';
import classnames from 'classnames';

export const Tooltip = ({children}) => <span className="tooltip">{children}</span>;

/**
 * Component for injecting tooltips into anything. Can be used in two ways:
 *
 * 1:
 * <WithTooltip>
 *   <button>
 *     Click me?
 *     <Tooltip>
 *       Click the button!
 *     </Tooltip>
 *   </button>
 * </WithTooltip>
 *
 * and 2:
 * <WithTooltip text="Click the button!">
 *   <button>
 *     Click me?
 *   </button>
 * </WithTooltip>
 *
 * Both of these examples will produce the same effect. The first method is
 * more generic and allows us to separate the hover area from the tooltip
 * location.
 */
export const WithTooltip = props => {
  const {children, shouldShow = true, generateText} = props;
  const text = generateText ? generateText(props) : props.text;
  const onlyChild = React.Children.only(children);
  const nextProps = _.assign({}, _.omit(props, 'children', 'shouldShow', 'generateText', 'text'), onlyChild.props, {className: classnames(onlyChild.props.className, 'withtooltip')});

  // If there is no text, then assume the tooltip is already nested.
  const tooltip = typeof text === 'string' && shouldShow && <Tooltip>{text}</Tooltip>;
  const nestedChildren = nextProps.children ? [].concat(nextProps.children, tooltip) : [tooltip];

  return React.cloneElement(onlyChild, nextProps, ...nestedChildren);
};
