new Vue({
    el: '#app',
    delimiters: ['${', '}'],
    vuetify: new Vuetify(),
    data() {
        return {
            keys: [[1, 2, 3, 'C'], [4, 5, 6], [7, 8, 9]],
            table: this.newTable(null),
            editable: this.newTable(true),
            selectedRow: null,
            selectedCol: null
        }
    },
    methods: {
        newTable: function(x) {
            var t = new Array(9)
            for (let i = 0; i < 9; i++) {
                t[i] = new Array(9).fill(x)
            }
            return t
        },
        selectCell: function(row, col) {
            this.selectedRow = row
            this.selectedCol = col
            return
        },
        isFocusedRow: function(row) {
            return row == this.selectedRow
        },
        isFocusedCol: function(col) {
            return col == this.selectedCol
        },
        isFocusedCell: function(row, col) {
            return this.isFocusedRow(row) && this.isFocusedCol(col)
        },
        setNumber: function(row, col, num) {
            if (0 <= row && row < 9 && 0 <= col && col < 9) {
                if (this.editable[row][col]) {
                    this.table[row].splice(col, 1, num)
                }
            }
        },
        keyAction: function(key) {
            let row = this.selectedRow
            let col = this.selectedCol
            if (key === 'C') {
                this.setNumber(row, col, null)
            } else {
                this.setNumber(row, col, key)
            }
        }
    },
    mounted: function() {
        for (let i = 0; i < 9; i++) {
            this.table[i].splice(i, 1, i+1)
            this.editable[i].splice(i, 1, false)
        }
    }
})