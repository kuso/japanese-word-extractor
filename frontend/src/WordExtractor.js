import React, { useState } from 'react'
import { usePostCallback } from "use-axios-react";
import Box from '@material-ui/core/Box';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import Tooltip from '@material-ui/core/Tooltip';
import { withStyles } from '@material-ui/core/styles';

/*
const HtmlTooltip = (theme => ({
    tooltip: {
        fontFamily: 'Arial,sans-serif !important',
        backgroundColor: '#f5f5f9',
        color: 'rgba(0, 0, 0, 0.87)',
        maxWidth: 220,
        //fontSize: theme.typography.pxToRem(12),
        border: '1px solid #dadde9',
    },
}))(Tooltip);

 */
const HtmlTooltip = withStyles(theme => ({
    tooltip: {
        fontFamily: 'Arial, sans-serif !important',
        backgroundColor: '#f5f5f9',
        color: 'rgba(0, 0, 0, 0.87)',
        maxWidth: 320,
        fontSize: theme.typography.pxToRem(28),
        border: '1px solid #dadde9',
    },
}))(Tooltip);

function useInput(initialValue) {
    const [value,setValue] = useState(initialValue);

    function handleChange(e){
        setValue(e.target.value);
    }

    return [value, handleChange];
}

const NormalToken = (props) => {
    const [color, setColor] = useState("white")
    const handleMouseOver = () => {
        setColor("blue")
    }

    const handleMouseOut = () => {
        setColor("white")
    }

    var token = props.token
    return (
        <HtmlTooltip
            title={
                <React.Fragment>
                    { token.DictFormHiragana }{' '}<br/>
                    { token.DictForm }{' '}<br/>
                    { token.Meaning }{' '}<br/>
                </React.Fragment>
            }
        >
            <span className={"vocab"} onMouseOver={handleMouseOver} onMouseOut={handleMouseOut}><ruby><rb>{token.Text}</rb><rt className={color}>{ token.DictFormHiragana }</rt></ruby></span>
        </HtmlTooltip>
    )
}

const JLPTToken = (props) => {
    var token = props.token
    const [color, setColor] = useState("white")

    const handleMouseOver = () => {
        setColor("blue")
    }

    const handleMouseOut = () => {
        setColor("white")
    }

    return (
        <HtmlTooltip
            title={
                <React.Fragment>
                    <Typography color="inherit"></Typography>
                    N{ token.Level}{' '}<br/>
                    { token.DictFormHiragana }{' '}<br/>
                    { token.DictForm }{' '}<br/>
                    { token.Meaning }{' '}<br/>
                </React.Fragment>
            }
        >
            <span className={"jlptn"+token.Level} onMouseOver={handleMouseOver} onMouseOut={handleMouseOut}><ruby><rb>{token.Text}</rb><rt className={color}>{ token.DictFormHiragana }</rt></ruby></span>
        </HtmlTooltip>
    )
}

function div(inner) {
    return (
        <div className="mb-3">
            {inner}
        </div>
    )
}

const ConvertedText = (data) => {
    var sections = data.sections
    var out = []
    var count = 1
    if (!sections) {
        return (
            <span>Please enter some Japanese text...</span>
        )
    }
    for (var i=0; i<sections.length; i++) {
        var section = sections[i]
        var sectionOut = []
        for (var j=0; j<section.tokens.length; j++) {
            count = count + 1
            var token = section.tokens[j]
            if (token.Level > 0) {
                sectionOut.push(<JLPTToken key={count} token={token}/>)
            } else {
                // case not in JLPT vocabs
                if (token.Text !== token.DictFormHiragana && token.DictFormHiragana !== "") {
                    sectionOut.push(<NormalToken key={count} token={token}/>)
                } else {
                    sectionOut.push(token.Text)
                }
            }
        }
        out.push(div(sectionOut))
    }
    return (
        <Box className="converted">{out}</Box>
    )
}

const WordExtractor = () => {
    const [queryText, setQueryText] = useInput("");
    const query = {id: "", querytext: queryText}

    function postQueryRequest({ id, querytext }) {
        return {
            url: "http://127.0.0.1:8081/v1/job",
            data: {id, querytext}
        };
    }

    /*
            {error && <Box><code>{JSON.stringify(error)}</code></Box>}
     */

    const StatusBar = ({loading, error}) => (
        <span>
            {loading && <span>Loading...</span>}
            {error && <span> Error sending requests...</span>}
        </span>
    );

    const [exec, loading, {error, data}] = usePostCallback(postQueryRequest);
    return (
        <Box my={3}>
            <Typography variant="h3" component="h3" gutterBottom>
                jwordlist.com
            </Typography>
            <Typography variant="h5" component="h5" gutterBottom>
                Extract Japanese Vocabulary From Text
            </Typography>
            <Box my={3}>
                <TextField variant="filled" label="Enter some Japanese text" multiline rows="10" fullWidth value={queryText} onChange={setQueryText}/>
            </Box>
            <Box my={3}>
                <Button variant="contained" onClick={() => exec(query)} color="primary">Extract Vocabulary</Button>
                <StatusBar loading={loading} error={error}/>
            </Box>
            {data && <ConvertedText sections={data.sections} />}
        </Box>
    )
}

export default WordExtractor