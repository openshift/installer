import React from 'react';

export const Loader = () => {
  return (
    <div className="wiz-wizard">
      <div className="wiz-wizard__cell wiz-wizard__content cos-loader-container">
        <div className="cos-loader">
          <div className="cos-loader-dot__one"></div>
          <div className="cos-loader-dot__two"></div>
          <div className="cos-loader-dot__three"></div>
        </div>
      </div>
    </div>
  );
};

export const LoaderInline = () => <div className="cos-loader-inline">
  <div className="cos-loader-dot__one"></div>
  <div className="cos-loader-dot__two"></div>
  <div className="cos-loader-dot__three"></div>
</div>;
