import React, { useState, useEffect, useCallback } from 'react'
import {
  GridList,
  GridListTile,
  Typography,
  Collapse,
  Link,
} from '@material-ui/core'
import { makeStyles } from '@material-ui/core/styles'
import withWidth from '@material-ui/core/withWidth'

import subsonic from '../subsonic'

import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardMedia from '@material-ui/core/CardMedia'
import PropTypes from 'prop-types'

import { AlbumGridTile } from '../album/AlbumGridView'
import { getColsForWidth } from '../album/AlbumGridView'
import {
  useTranslate,
  useShowController,
  ShowContextProvider,
  useQueryWithStore,
  Loading,
} from 'react-admin'
import { Redirect } from 'react-router'

const useStyles = makeStyles(
  (theme) => ({
    root: {
      display: 'flex',
      padding: '1em',
      '& .MuiTypography-h5': {
        wordBreak: 'break-word',
      },
      [theme.breakpoints.down('xs')]: {
        padding: 'unset',
        background: ({ img }) => `url(${img})`,
      },
    },
    bgContainer: {
      display: 'flex',
      width: '100%',
      [theme.breakpoints.down('xs')]: {
        height: '15rem',
        width: '100vw',
        padding: 'unset',
        backdropFilter: 'blur(1px)',
        backgroundPosition: '50% 30%',
        background: `linear-gradient(to bottom, rgba(52 52 52 / 72%), rgba(21 21 21))`,
      },
    },
    iroot: {
      margin: '20px',
      display: 'grid',
    },
    details: {
      display: 'flex',
      flex: '1',
      flexDirection: 'column',
    },
    link: {
      margin: '1px',
    },
    mdetails: {
      display: 'none',
      [theme.breakpoints.down('xs')]: {
        display: 'flex',
        alignItems: 'center',
        width: '7rem',
        marginLeft: '0.5rem',
        flex: '1',
      },
    },
    mbio: {
      display: 'none',
      [theme.breakpoints.down('xs')]: {
        display: 'flex',
        marginLeft: '3%',
        marginRight: '3%',
        zIndex: '1',
        '& p': {
          whiteSpace: ({ expanded }) => (expanded ? 'unset' : 'nowrap'),
          overflow: 'hidden',
          width: '95vw',
          textOverflow: 'ellipsis',
        },
      },
    },
    content: {
      flex: '1 0 auto',
      '& .MuiTypography-root': {
        display: ({ expanded }) => (expanded ? 'block' : '-webkit-inline-box'),
        boxOrient: 'vertical',
        lineClamp: '3',
      },
    },
    cover: {
      width: 151,
      boxShadow: '0px 0px 6px 0px #565656',
      borderRadius: '5px',
      [theme.breakpoints.up('sm')]: {
        borderRadius: '7em',
      },
    },
    martImage: {
      marginLeft: '1em',
      maxHeight: '10rem',
      backgroundColor: 'inherit',
      display: 'none',
      [theme.breakpoints.down('xs')]: {
        marginTop: '4rem',
        maxHeight: '7rem',
        width: '7rem',
        display: 'flex',
      },
    },
    artImage: {
      maxHeight: '9.5rem',
      backgroundColor: 'inherit',
      display: 'flex',
      [theme.breakpoints.down('xs')]: {
        marginTop: '4rem',
        maxHeight: '7rem',
        width: '7rem',
      },
    },
    artDetail: {
      flex: '1',
      padding: '3%',
      display: 'flex',
      minHeight: '10rem',
      '& .MuiPaper-elevation1': {
        boxShadow: 'none',
        padding: '4px',
      },
      [theme.breakpoints.down('xs')]: {
        display: 'none',
      },
    },
    album: {
      marginBottom: '1em',
    },
  }),
  { name: 'NDArtistPage' }
)

const ImgMediaCard = ({ artistId, artist }) => {
  const [artistInfo, setartistInfo] = useState()
  const [expanded, setExpanded] = useState(false)

  const title = artist.artist
  let completeBioLink = ''
  const link = artistInfo?.biography?.match(
    /<a\s+(?:[^>]*?\s+)?href=(["'])(.*?)\1/
  )
  const biography = artistInfo?.biography?.replace(new RegExp('<.*>', 'g'), '')
  const translate = useTranslate()

  const handleExpandClick = useCallback(() => {
    setExpanded(!expanded)
  }, [expanded, setExpanded])

  if (link) {
    completeBioLink = link[2]
  }

  useEffect(() => {
    subsonic
      .getArtistInfo(artistId)
      .then((resp) => resp.json['subsonic-response'])
      .then((data) => {
        if (data.status === 'ok') {
          setartistInfo(data.artistInfo)
        }
      })
      .catch((e) => {
        console.error('error on artist page', e)
        return (
          <Redirect
            to={`/album?filter={"artist_id":"${artistId}"}&order=ASC&sort=maxYear&
              displayedFilters={"compilation":true}`}
          />
        )
      })
  }, [artistId, artist])

  const img = artistInfo?.largeImageUrl
  const classes = useStyles({ img, link, expanded })

  return (
    <>
      <div className={classes.root}>
        <div className={classes.bgContainer}>
          <Card className={classes.martImage}>
            <CardMedia
              className={classes.cover}
              image={`${artistInfo?.mediumImageUrl}`}
              title={title}
            />
          </Card>
          <div className={classes.mdetails}>
            <Typography component="h5" variant="h5">
              {title}
            </Typography>
          </div>
          <Card className={classes.artDetail}>
            <Card className={classes.artImage}>
              <CardMedia
                className={classes.cover}
                image={`${artistInfo?.mediumImageUrl}`}
                title={title}
              />
            </Card>
            <div className={classes.details}>
              <CardContent className={classes.content}>
                <Typography component="h5" variant="h5">
                  {title}
                </Typography>
                <Collapse
                  collapsedHeight={'4.5em'}
                  in={expanded}
                  timeout={'auto'}
                >
                  <Typography variant={'body1'} onClick={handleExpandClick}>
                    {biography}
                    <Link
                      href={completeBioLink}
                      className={classes.link}
                      target="_blank"
                      rel="nofollow"
                    >
                      {translate('message.lastfmLink')}
                    </Link>
                  </Typography>
                </Collapse>
              </CardContent>
            </div>
          </Card>
        </div>
      </div>
      <div className={classes.mbio}>
        <Collapse collapsedHeight={'1.5em'} in={expanded} timeout={'auto'}>
          <Typography variant={'body1'} onClick={handleExpandClick}>
            {biography}
            <Link
              href={completeBioLink}
              className={classes.link}
              target="_blank"
              rel="nofollow"
            >
              {translate('message.lastfmLink')}
            </Link>
          </Typography>
        </Collapse>
      </div>
    </>
  )
}

const ArtistAlbum = ({ artist, record, width }) => {
  const classes = useStyles()
  const translate = useTranslate()

  return (
    <>
      <ImgMediaCard artistId={record.id} artist={artist[0]} />
      <div className={classes.iroot}>
        <div className={classes.album}>
          {artist.length +
            ' ' +
            translate('resources.album.name', { smart_count: artist.length })}
        </div>
        <GridList
          component={'div'}
          cellHeight={'auto'}
          cols={getColsForWidth(width)}
          spacing={20}
        >
          {artist.map((artist) => (
            <GridListTile className={classes.gridListTile} key={artist.id}>
              <AlbumGridTile
                record={artist}
                basePath={'/album'}
                showArtist={true}
              />
            </GridListTile>
          ))}
        </GridList>
      </div>
    </>
  )
}

const ArtistShow = (props) => {
  const { width } = props

  const controllerProps = useShowController(props)
  const { record } = { ...controllerProps }

  const payload = {
    pagination: { page: 1, perPage: 12 },
    sort: { field: 'name', order: 'ASC' },
    filter: { artist_id: record?.id },
  }
  const { loaded, data } = useQueryWithStore({
    type: 'getList',
    resource: 'album',
    payload: payload,
  })
  if (!loaded) {
    return <Loading />
  }

  return (
    <ShowContextProvider value={controllerProps}>
      <ArtistAlbum
        width={width}
        artist={data}
        {...props}
        {...controllerProps}
      />
    </ShowContextProvider>
  )
}

ArtistShow.propTypes = {
  width: PropTypes.oneOf(['lg', 'md', 'sm', 'xl', 'xs']),
}

export default withWidth()(ArtistShow)
